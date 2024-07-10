package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-relyt/internal/provider/client"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dpsResource{}
	_ resource.ResourceWithConfigure   = &dpsResource{}
	_ resource.ResourceWithImportState = &dpsResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDpsResource() resource.Resource {
	return &dpsResource{}
}

// orderResource is the resource implementation.
type dpsResource struct {
	client *client.RelytClient
}

// Metadata returns the resource type name.
func (r *dpsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dps"
}

// Schema defines the schema for the resource.
func (r *dpsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dwsu":        schema.ObjectAttribute{Computed: true},
			"name":        schema.StringAttribute{Required: true},
			"description": schema.StringAttribute{Required: true},
			"engine":      schema.StringAttribute{Required: true},
			"size":        schema.Int64Attribute{Required: true},
		},
	}
}

// Create a new resource.
func (r *dpsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from dpsModel
	var dpsModel DpsModel
	diags := req.Plan.Get(ctx, &dpsModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	relytDps := client.DpsMode{
		Description: dpsModel.Description.ValueString(),
		Engine:      dpsModel.Engine.ValueString(),
		Name:        dpsModel.Name.ValueString(),
		Spec: client.Spec{
			ID: dpsModel.Size.ValueInt64(),
		},
	}

	// Create new order
	createResult, err := r.client.CreateEdps(ctx, dpsModel.Dwsu.ID.ValueString(), relytDps)
	if err != nil || createResult.Code != 200 {
		resp.Diagnostics.AddError(
			"Error creating dps",
			"Could not create dps, unexpected error: "+err.Error(),
		)
		return
	}
	queryDpsMode, err := r.client.TimeOutTask(600, func() (any, error) {
		return r.client.GetDps(ctx, dpsModel.Dwsu.ID.ValueString(), createResult.Data)
	})
	if err != nil {
		tflog.Error(ctx, "error wait dps ready"+err.Error())
		resp.Diagnostics.AddError(
			"Error creating dps",
			"Could not create dps, unexpected error: "+err.Error(),
		)
		return
		//fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
	relytQueryModel := queryDpsMode.(*client.DpsMode)
	tflog.Info(ctx, "bizId:"+relytQueryModel.ID)
	// 将毫秒转换为秒和纳秒
	seconds := relytQueryModel.UpdateTime / 1000
	nanoseconds := (relytQueryModel.UpdateTime % 1000) * int64(time.Millisecond)

	// 使用 time.Unix 函数创建 time.Time 对象
	t := time.Unix(seconds, nanoseconds)
	dpsModel.LastUpdated = types.StringValue(t.Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, dpsModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *dpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DwsuModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	dwsu, err := r.client.GetDwsu(ctx, state.ID.ValueString())
	if err != nil {
		tflog.Error(ctx, "error read dwsu"+err.Error())
		return
	}
	state.Alias = types.StringValue(dwsu.Alias)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dpsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state client.DwsuModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DropDwsu(ctx, state.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Order",
			"Could not delete order, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *dpsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	relytClient, ok := req.ProviderData.(*client.RelytClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *RelytClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = relytClient
}

func (r *dpsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
