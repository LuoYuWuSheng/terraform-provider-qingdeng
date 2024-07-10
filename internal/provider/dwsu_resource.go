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
	_ resource.Resource                = &dwsuResource{}
	_ resource.ResourceWithConfigure   = &dwsuResource{}
	_ resource.ResourceWithImportState = &dwsuResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDwsuResource() resource.Resource {
	return &dwsuResource{}
}

// orderResource is the resource implementation.
type dwsuResource struct {
	client *client.RelytClient
}

// Metadata returns the resource type name.
func (r *dwsuResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dwsu"
}

// Schema defines the schema for the resource.
func (r *dwsuResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"cloud":     schema.StringAttribute{Required: true},
			"region":    schema.StringAttribute{Required: true},
			"dwsu_type": schema.StringAttribute{Required: true},
			"alias":     schema.StringAttribute{Optional: true},
			"defaultDps": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name":        schema.StringAttribute{Required: true},
					"description": schema.StringAttribute{Optional: true},
					"engine":      schema.StringAttribute{Required: true},
					"size":        schema.Int64Attribute{Required: true},
				},
			},
		},
	}
}

// Create a new resource.
func (r *dwsuResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from dwsuModel
	var dwsuModel DwsuModel
	diags := req.Plan.Get(ctx, &dwsuModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	relytDwsu := client.DwsuModel{
		DefaultDps: client.DpsMode{
			Description: dwsuModel.DefaultDps.Description.ValueString(),
			Engine:      dwsuModel.DefaultDps.Engine.ValueString(),
			Name:        dwsuModel.DefaultDps.Name.ValueString(),
			Spec: client.Spec{
				ID: dwsuModel.DefaultDps.Size.ValueInt64(),
			},
		},
		Region: client.Region{
			Cloud: client.Cloud{
				ID: dwsuModel.Cloud.ValueString(),
			},
			ID: dwsuModel.Region.ValueString(),
		},
	}

	// Create new order
	createResult, err := r.client.CeateDwsu(ctx, relytDwsu)
	if err != nil || createResult.Code != 200 {
		resp.Diagnostics.AddError(
			"Error creating dwsu",
			"Could not create dwsu, unexpected error: "+err.Error(),
		)
		return
	}
	queryDwsuModel, err := r.client.TimeOutTask(600, func() (any, error) {
		return r.client.GetDwsuByAlias(ctx, dwsuModel.Alias.ValueString())
	})
	if err != nil {
		tflog.Error(ctx, "error wait dwsu ready"+err.Error())
		return
		//fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
	relytQueryModel := queryDwsuModel.(*client.DwsuModel)
	tflog.Info(ctx, "bizId:"+relytQueryModel.ID)
	dwsuModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, dwsuModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *dwsuResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
func (r *dwsuResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dwsuResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
func (r *dwsuResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dwsuResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
