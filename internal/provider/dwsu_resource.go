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
			"domain":    schema.StringAttribute{Required: true},
			"variant":   schema.StringAttribute{Computed: true},
			"edition":   schema.StringAttribute{Computed: true},
			"region":    schema.StringAttribute{Required: true},
			"dwsu_type": schema.StringAttribute{Required: true},
			"alias":     schema.StringAttribute{Optional: true},
			//"last_updated": schema.Int64Attribute{Computed: true},
			//"status":       schema.StringAttribute{Computed: true},
			"default_dps": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name":        schema.StringAttribute{Required: true},
					"description": schema.StringAttribute{Optional: true},
					"engine":      schema.StringAttribute{Required: true},
					"size":        schema.StringAttribute{Required: true},
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
	dwsuModel.Variant = types.StringValue("basic")
	dwsuModel.Edition = types.StringValue("standard")
	if resp.Diagnostics.HasError() {
		return
	}
	relytDwsu := client.DwsuModel{
		DefaultDps: &client.DpsMode{
			Description: dwsuModel.DefaultDps.Description.ValueString(),
			Engine:      dwsuModel.DefaultDps.Engine.ValueString(),
			Name:        dwsuModel.DefaultDps.Name.ValueString(),
			Spec: &client.Spec{
				Name: dwsuModel.DefaultDps.Size.ValueString(),
			},
		},
		Domain:  dwsuModel.Domain.ValueString(),
		Alias:   dwsuModel.Alias.ValueString(),
		Variant: &client.Variant{ID: dwsuModel.Variant.ValueString()},
		Edition: &client.Edition{ID: dwsuModel.Edition.ValueString()},
		Region: &client.Region{
			Cloud: &client.Cloud{
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
	if createResult.Data == nil {
		resp.Diagnostics.AddError(
			"Error creating dwsu",
			"Could not get dwsu id, after create!",
		)
		return
	}
	dwsuModel.ID = types.StringValue(*createResult.Data)
	queryDwsuModel, err := r.client.TimeOutTask(600, func() (any, error) {
		dwsu, err2 := r.client.GetDwsu(ctx, *createResult.Data)
		if err2 != nil || dwsu == nil {
			//这里判断是否要充实
			return dwsu, err2
		}
		//todo 不判断defaultDps状态,读接口没返回这个状态
		if dwsu.Status == client.DWSU_STATUS_READY {
			return dwsu, nil
		}
		return dwsu, fmt.Errorf("dwsu is not Ready")
	})
	if err != nil {
		tflog.Error(ctx, "error wait dwsu ready"+err.Error())
		return
		//fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
	relytQueryModel := queryDwsuModel.(*client.DwsuModel)
	tflog.Info(ctx, "bizId:"+relytQueryModel.ID)
	//dwsuModel.LastUpdated = types.Int64Value(time.Now().UnixMilli())
	//dwsuModel.Status = types.StringValue(relytQueryModel.Status)
	// Set state to fully populated data
	diags = resp.State.Set(ctx, dwsuModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create dwsu success: "+relytQueryModel.ID)
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
	_, err := r.client.GetDwsu(ctx, state.ID.ValueString())
	if err != nil {
		tflog.Error(ctx, "error read dwsu"+err.Error())
		return
	}
	//state.Status = types.StringValue(dwsu.Status)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read dwsu succ : "+state.ID.ValueString())
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dwsuResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not support", "update dwsu not supported")
	return
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dwsuResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DwsuModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DropDwsu(ctx, state.ID.ValueString())
	if err != nil {
		//要不要加error
		resp.Diagnostics.AddWarning(
			"Error Deleting dwsu",
			"Could not delete dwsu, unexpected error: "+err.Error(),
		)
		//return
	}
	//等待删除完成
	_, err = r.client.TimeOutTask(600, func() (any, error) {
		dwsu, err2 := r.client.GetDwsu(ctx, state.ID.ValueString())
		if err2 != nil || dwsu == nil {
			//这里判断是否要充实
			return dwsu, err2
		}
		if dwsu == nil || dwsu.Status == client.DWSU_STATUS_DROPPED {
			return dwsu, nil
		}
		return dwsu, fmt.Errorf("wait delete dwsu timeout ")
	})
	return
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
