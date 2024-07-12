package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RelytProviderModel describes the provider data model.
type RelytProviderModel struct {
	ApiHost   types.String `tfsdk:"api_host"`
	AuthKey   types.String `tfsdk:"auth_key"`
	Role      types.String `tfsdk:"role"`
	RegionApi types.String `tfsdk:"region_api"`
}

//todo !!!!! 这里的格式增减一定要注意！新增的属性如果老字段不存在，terraform会认为该资源被污染，会走删除重建流程

type DefaultDps struct {
	Description types.String `tfsdk:"description"`
	Engine      types.String `tfsdk:"engine"`
	Name        types.String `tfsdk:"name"`
	Size        types.String `tfsdk:"size"`
}

type DwsuModel struct {
	ID         types.String `tfsdk:"id"`
	Alias      types.String `tfsdk:"alias"`
	Cloud      types.String `tfsdk:"cloud"`
	Region     types.String `tfsdk:"region"`
	DwsuType   types.String `tfsdk:"dwsu_type"`
	Variant    types.String `tfsdk:"variant"`
	Edition    types.String `tfsdk:"edition"`
	Domain     types.String `tfsdk:"domain"`
	DefaultDps DefaultDps   `tfsdk:"default_dps"`
	//LastUpdated types.Int64  `tfsdk:"last_updated"`
	//Status      types.String `tfsdk:"status"`
}

type DpsModel struct {
	DwsuId      types.String `tfsdk:"dwsu_id"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Engine      types.String `tfsdk:"engine"`
	Size        types.String `tfsdk:"size"`
	//LastUpdated types.String `tfsdk:"last_updated"`
	//Status      types.String `tfsdk:"status"`
}

type DWUserModel struct {
	DwsuId                             types.String `tfsdk:"dwsu_id"`
	ID                                 types.String `tfsdk:"id"`
	AccountName                        types.String `tfsdk:"account_name"`
	AccountPassword                    types.String `tfsdk:"account_password"`
	DatalakeAwsLakeformationRoleArn    types.String `tfsdk:"datalake_aws_lakeformation_role_arn"`
	AsyncQueryResultLocationPrefix     types.String `tfsdk:"async_query_result_location_prefix"`
	AsyncQueryResultLocationAwsRoleArn types.String `tfsdk:"async_query_result_location_aws_role_arn"`
	//LastUpdated                        types.String `tfsdk:"last_updated"`
	//Status                             types.String `tfsdk:"status"`
}
