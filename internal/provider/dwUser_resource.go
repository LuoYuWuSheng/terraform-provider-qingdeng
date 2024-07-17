package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-relyt/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dwUserResource{}
	_ resource.ResourceWithConfigure   = &dwUserResource{}
	_ resource.ResourceWithImportState = &dwUserResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewdwUserResource() resource.Resource {
	return &dwUserResource{}
}

// orderResource is the resource implementation.
type dwUserResource struct {
	client *client.RelytClient
}

// Metadata returns the resource type name.
func (r *dwUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dwuser"
}

// Schema defines the schema for the resource.
func (r *dwUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"dwsu_id": schema.StringAttribute{Required: true, Description: "The ID of the service unit."},
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "The ID of the DW user.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"account_name":                             schema.StringAttribute{Required: true, Description: "The name of the DW user, which is unique in the instance. The name is the email address."},
			"account_password":                         schema.StringAttribute{Required: true, Description: "initPassword"},
			"datalake_aws_lakeformation_role_arn":      schema.StringAttribute{Optional: true, Description: ""},
			"async_query_result_location_prefix":       schema.StringAttribute{Optional: true},
			"async_query_result_location_aws_role_arn": schema.StringAttribute{Optional: true},
		},
	}
}

// Create a new resource.
func (r *dwUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from dwUserModel
	var dwUserModel DWUserModel
	diags := req.Plan.Get(ctx, &dwUserModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	meta := RouteRegionUri(ctx, dwUserModel.DwsuId.ValueString(), r.client, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	regionUri := meta.URI
	relytAccount := client.Account{
		InitPassword: dwUserModel.AccountPassword.ValueString(),
		Name:         dwUserModel.AccountName.ValueString(),
	}
	asyncResult := client.AsyncResult{
		AwsIamArn:        dwUserModel.AsyncQueryResultLocationAwsRoleArn.ValueString(),
		S3LocationPrefix: dwUserModel.AsyncQueryResultLocationPrefix.ValueString(),
	}
	lakeFormation := client.LakeFormation{
		IAMRole: dwUserModel.DatalakeAwsLakeformationRoleArn.ValueString(),
	}

	// Create new order
	createResult, err := r.client.CreateAccount(ctx, regionUri, dwUserModel.DwsuId.ValueString(), relytAccount)
	if err != nil || createResult.Code != 200 {
		resp.Diagnostics.AddError(
			"Error creating dwuser",
			"Could not create dwuser, unexpected error: "+err.Error(),
		)
		return
	}
	//dwUserModel.ID = types.StringValue(*createResult.Data)
	//tflog.Info(ctx, "accountId:"+*createResult.Data)
	dwUserModel.ID = dwUserModel.AccountName
	_, err = r.client.AsyncAccountConfig(ctx, regionUri, dwUserModel.DwsuId.ValueString(), dwUserModel.ID.ValueString(), asyncResult)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error config dwuser",
			"Could not config dwuser async, unexpected error: "+err.Error(),
		)
		return
	}
	_, err = r.client.LakeFormationConfig(ctx, regionUri, dwUserModel.DwsuId.ValueString(), dwUserModel.ID.ValueString(), lakeFormation)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error config dwuser",
			"Could not config dwuser lakeformation, unexpected error: "+err.Error(),
		)
		return
	}
	// 将毫秒转换为秒和纳秒
	//seconds := relytQueryModel.UpdateTime / 1000
	//nanoseconds := (relytQueryModel.UpdateTime % 1000) * int64(time.Millisecond)

	// 使用 time.Unix 函数创建 time.Time 对象
	//t := time.Unix(seconds, nanoseconds)
	//dwUserModel.LastUpdated = types.StringValue(t.Format(time.RFC850))
	// Set state to fully populated data
	diags = resp.State.Set(ctx, dwUserModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *dwUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DWUserModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//meta := RouteRegionUri(ctx, state.DwsuId.ValueString(), r.client, &resp.Diagnostics)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//regionUri := meta.URI
	//relytAccount := client.Account{
	//	InitPassword: state.AccountPassword.ValueString(),
	//	Name:         state.AccountName.ValueString(),
	//}
	//dwuser, err := r.client.CreateAccount(ctx, "", state.DwsuId.ValueString(), relytAccount)
	//if err != nil {
	//	tflog.Error(ctx, "error read dwuer"+err.Error())
	//	return
	//}
	//state.ID = types.StringValue(*dwuser.Data)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dwUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//resp.Diagnostics.AddError("not support", "update account not supported")
	var state DWUserModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.ID = state.AccountName
	meta := RouteRegionUri(ctx, state.DwsuId.ValueString(), r.client, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	regionUri := meta.URI
	tflog.Info(ctx, "accountId:"+state.ID.ValueString())
	asyncResult := client.AsyncResult{
		AwsIamArn:        state.AsyncQueryResultLocationAwsRoleArn.ValueString(),
		S3LocationPrefix: state.AsyncQueryResultLocationPrefix.ValueString(),
	}
	lakeFormation := client.LakeFormation{
		IAMRole: state.DatalakeAwsLakeformationRoleArn.ValueString(),
	}
	_, err := r.client.AsyncAccountConfig(ctx, regionUri, state.DwsuId.ValueString(), state.ID.ValueString(), asyncResult)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error config dwuser",
			"Could not config dwuser async, unexpected error: "+err.Error(),
		)
		return
	}
	_, err = r.client.LakeFormationConfig(ctx, regionUri, state.DwsuId.ValueString(), state.ID.ValueString(), lakeFormation)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error config dwuser",
			"Could not update dwuser lakeformation, unexpected error: "+err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	return
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dwUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DWUserModel
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

	// Delete existing order
	err := r.client.DropAccount(ctx, regionUri, state.DwsuId.ValueString(), state.ID.ValueString())
	if err != nil {
		//要不要加error
		resp.Diagnostics.AddError(
			"Error Deleting dwuser",
			"Could not delete dwuser, unexpected error: "+err.Error(),
		)
	}
}

// Configure adds the provider configured client to the resource.
func (r *dwUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	relytClient, ok := req.ProviderData.(*client.RelytClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *RelytClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = relytClient
}

func (r *dwUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
