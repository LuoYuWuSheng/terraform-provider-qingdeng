// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
)

// Ensure RelytProvider satisfies various provider interfaces.
var _ provider.Provider = &RelytProvider{}
var _ provider.ProviderWithFunctions = &RelytProvider{}

// RelytProvider defines the provider implementation.
type RelytProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RelytProviderModel describes the provider data model.
type RelytProviderModel struct {
	ApiHost types.String `tfsdk:"api_host"`
	AuthKey types.String `tfsdk:"auth_key"`
	RoleId  types.String `tfsdk:"role_id"`
}

type RelytProviderEnv struct {
	EnvKey         string
	PropertyName   string
	SummarySuggest string
	detailSuggest  string
}

var (
	apiHostEnv = RelytProviderEnv{
		EnvKey:         "RELYT_API_HOST",
		PropertyName:   "api_host",
		SummarySuggest: "Unknown Relyt API Host",
		detailSuggest: "The provider cannot create the Relyt API client as there is an unknown configuration value for the Relyt API apiHost. " +
			"Either target apply the source of the value first, set the value statically in the configuration, or use the RELYT_API_HOST environment variable.",
	}
	authKeyEnv = RelytProviderEnv{
		EnvKey:         "RELYT_AUTH_KEY",
		PropertyName:   "auth_key",
		SummarySuggest: "Unknown Relyt Auth EnvKey",
		detailSuggest: "The provider cannot create the Relyt API client as there is an unknown configuration value for the Relyt Auth EnvKey. " +
			"Either target apply the source of the value first, set the value statically in the configuration, or use the RELYT_AUTH_KEY environment variable.",
	}
	roleIDEnv = RelytProviderEnv{
		EnvKey:         "RELYT_ROLE_ID",
		PropertyName:   "role_id",
		SummarySuggest: "Unknown Relyt Role ID",
		detailSuggest: "The provider cannot create the Relyt API client as there is an unknown configuration value for the Relyt Relyt Role ID. " +
			"Either target apply the source of the value first, set the value statically in the configuration, or use the RELYT_ROLE_ID environment variable.",
	}
)

func (p *RelytProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "qingdeng"
	resp.Version = p.version
}

// 定义provider能接受的参数，类型，是否可选等
func (p *RelytProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	tflog.Info(ctx, "===== scheme get ")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			//"endpoint": schema.StringAttribute{
			//	MarkdownDescription: "Example provider attribute",
			//	Optional:            true,
			//},
			"api_host": schema.StringAttribute{
				Required: true,
			},
			"auth_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"role_id": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// 读取配置文件
func (p *RelytProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data RelytProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if data.ApiHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(apiHostEnv.PropertyName),
			apiHostEnv.SummarySuggest,
			apiHostEnv.detailSuggest,
		)
	}

	if data.AuthKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(authKeyEnv.PropertyName),
			authKeyEnv.SummarySuggest,
			authKeyEnv.detailSuggest,
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	apiHost := os.Getenv(apiHostEnv.EnvKey)
	if !data.ApiHost.IsNull() {
		apiHost = data.ApiHost.ValueString()
	}
	authKey := os.Getenv(authKeyEnv.EnvKey)
	if !data.ApiHost.IsNull() {
		authKey = data.AuthKey.ValueString()
	}
	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiHost == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root(apiHostEnv.PropertyName),
			"Missing Relyt API Host",
			"The provider cannot create the Relyt API client as there is a missing or empty value for the Relyt API apiHost. "+
				"Set the apiHost value in the configuration or use the "+apiHostEnv.EnvKey+" environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if authKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root(authKeyEnv.PropertyName),
			"Missing Relyt AUTH KEY",
			"The provider cannot create the Relyt API client as there is a missing or empty value for the Relyt AUTH KEY. "+
				"Set the apiHost value in the configuration or use the "+authKeyEnv.EnvKey+" environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Example client configuration for data sources and resources

	// Create a new Relyt client using the configuration values
	tflog.Info(ctx, "host:"+apiHost+" auth:"+authKey+"role_id"+data.RoleId.ValueString())
	roleId := data.RoleId.ValueString()
	hashiCupsClient, err := hashicups.NewClient(&apiHost, &authKey, &roleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Relyt API Client",
			"An unexpected error occurred when creating the Relyt API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Relyt Client Error: "+err.Error(),
		)
		return
	}
	//client := http.DefaultClient
	client := hashiCupsClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *RelytProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrderResource,
		NewDpsResource,
	}
}

func (p *RelytProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	tflog.Info(ctx, "===== datasource get ")
	return []func() datasource.DataSource{
		//NewExampleDataSource,
		NewCoffeesDataSource,
	}
}

func (p *RelytProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RelytProvider{
			version: version,
		}
	}
}
