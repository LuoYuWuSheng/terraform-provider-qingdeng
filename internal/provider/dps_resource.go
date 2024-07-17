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
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"dwsu_id":     schema.StringAttribute{Required: true, Description: "The ID of the service unit."},
			"name":        schema.StringAttribute{Required: true, Description: "The name of the DPS cluster."},
			"id":          schema.StringAttribute{Computed: true, Description: "The ID of the DPS cluster."},
			"description": schema.StringAttribute{Optional: true, Description: "The description of the DPS cluster."},
			"engine":      schema.StringAttribute{Required: true, Description: "The type of the DPS cluster. enum:{extreme}"},
			"size":        schema.StringAttribute{Required: true, Description: "The name of the DPS cluster specification."},
			//"last_updated": schema.StringAttribute{Computed: true},
			//"status":       schema.StringAttribute{Computed: true},
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
	meta := RouteRegionUri(ctx, dpsModel.DwsuId.ValueString(), r.client, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	regionUri := meta.URI
	relytDps := client.DpsMode{
		Description: dpsModel.Description.ValueString(),
		Engine:      dpsModel.Engine.ValueString(),
		Name:        dpsModel.Name.ValueString(),
		Spec: &client.Spec{
			Name: dpsModel.Size.ValueString(),
		},
	}
	if dpsModel.ID.IsUnknown() {
		// Create new dps
		createResult, err := r.client.CreateEdps(ctx, regionUri, dpsModel.DwsuId.ValueString(), relytDps)
		if err != nil || createResult.Code != 200 {
			resp.Diagnostics.AddError(
				"Error creating dps",
				"Could not create dps, unexpected error: "+err.Error(),
			)
			return
		}
		if createResult.Data == nil {
			resp.Diagnostics.AddError(
				"Error creating dps",
				"Could not get dps id, after create!",
			)
			return
		}
		dpsModel.ID = types.StringValue(*createResult.Data)
		//拿到ID先写入，保障可重入。除非panic
		diags = resp.State.Set(ctx, dpsModel)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	queryDpsMode, err := r.client.TimeOutTask(r.client.CheckTimeOut, r.client.CheckInterval, func() (any, error) {
		dps, err2 := r.client.GetDps(ctx, regionUri, dpsModel.DwsuId.ValueString(), dpsModel.ID.ValueString())
		if err2 != nil {
			//这里判断是否要充实
			return dps, err2
		}
		if dps != nil && dps.Status == client.DPS_STATUS_READY {
			return dps, nil
		}
		return dps, fmt.Errorf("dps is not Ready")
	})
	if err != nil {
		tflog.Error(ctx, "error wait dps ready"+err.Error())
		resp.Diagnostics.AddError(
			"create dps failed!", "error wait dps ready! "+err.Error(),
		)
		return
		//fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
	relytQueryModel := queryDpsMode.(*client.DpsMode)
	tflog.Info(ctx, "bizId:"+relytQueryModel.ID)
	// 将毫秒转换为秒和纳秒
	//seconds := relytQueryModel.UpdateTime / 1000
	//nanoseconds := (relytQueryModel.UpdateTime % 1000) * int64(time.Millisecond)

	// 使用 time.Unix 函数创建 time.Time 对象
	//t := time.Unix(seconds, nanoseconds)
	//dpsModel.LastUpdated = types.StringValue(t.Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, dpsModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create dps succ !"+relytQueryModel.ID)
}

// Read resource information.
func (r *dpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DpsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	meta := RouteRegionUri(ctx, state.DwsuId.ValueString(), r.client, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	regionUri := meta.URI
	_, err := r.client.GetDps(ctx, regionUri, state.DwsuId.ValueString(), state.ID.ValueString())
	if err != nil {
		tflog.Error(ctx, "error read dps"+err.Error())
		return
	}
	//state.Status = types.StringValue(dps.Status)
	// Set refreshed state
	//尝试修改其中某些属性，看terraform行为
	//state.Description = types.StringValue("change desc")
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dpsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not support", "update dps not supported")
	return
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DpsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	meta := RouteRegionUri(ctx, state.DwsuId.ValueString(), r.client, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	regionUri := meta.URI

	// fixme 这里只要我删除，再读立刻为null了。。这样无法等待edps删除完成
	err := r.client.DropEdps(ctx, regionUri, state.DwsuId.ValueString(), state.ID.ValueString())
	if err != nil {
		tflog.Error(ctx, "error delete dps "+err.Error())
		resp.Diagnostics.AddError(
			"Error Deleting dps ",
			"Could not delete dps, unexpected error: "+err.Error(),
		)
		return
	}
	_, err = r.client.TimeOutTask(r.client.CheckTimeOut, r.client.CheckInterval, func() (any, error) {
		dps, err2 := r.client.GetDps(ctx, regionUri, state.DwsuId.ValueString(), state.ID.ValueString())
		if err2 != nil {
			//这里判断是否要重试
			return dps, err2
		}
		if dps == nil || dps.Status == client.DPS_STATUS_DROPPED {
			return dps, nil
		}
		return dps, fmt.Errorf("wait delete dps timeout ")
	})
	if err != nil {
		tflog.Error(ctx, "error wait dps delete "+err.Error())
		resp.Diagnostics.AddError(
			"Error Deleting dps ",
			"Could not delete dps, unexpected error: "+err.Error(),
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
