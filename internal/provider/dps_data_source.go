package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-relyt/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &dpsDataSource{}
	_ datasource.DataSourceWithConfigure = &dpsDataSource{}
)

// coffeesModel maps coffees schema data.
type SpecQueryModel struct {
	//standard
	Edition types.String `tfsdk:"edition"`
	//hybrid extreme
	Type     types.String `tfsdk:"type"`
	Cloud    types.String `tfsdk:"cloud"`
	Region   types.String `tfsdk:"region"`
	SpecName types.String `tfsdk:"spec_name"`
}

type SpecModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewDpsDataSource() datasource.DataSource {
	return &dpsDataSource{}
}

type dpsDataSource struct {
	client *client.RelytClient
}

func (d *dpsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_spec"
}

// Schema defines the schema for the data source.
func (d *dpsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"type":      schema.StringAttribute{Required: true},
			"edition":   schema.StringAttribute{Required: true},
			"spec_name": schema.StringAttribute{Required: true},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dpsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SpecQueryModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.Edition.IsNull() {
		state.Edition = types.StringValue("standard")
	}
	switch state.Type.ValueString() {
	case "hybrid":
	case "extreme":
		break
	default:
		resp.Diagnostics.AddError("data config err", "not support dpsType! only support hybrid or extreme")
		return
	}

	specs, err := d.client.ListSpec(ctx, state.Edition.ValueString(), state.Type.ValueString(), state.Cloud.ValueString(), state.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read relyt spec list api",
			err.Error(),
		)
		return
	}
	var target *SpecModel
	if len(specs) > 0 {
		for _, value := range specs {
			if state.SpecName.ValueString() == value.Name {
				target = &SpecModel{
					ID:   types.Int64Value(value.ID),
					Name: types.StringValue(value.Name),
				}
			}
		}
	}
	if target == nil {
		resp.Diagnostics.AddError("spec not find", "can't find match spec "+state.SpecName.ValueString())
	}

	// Set state
	diags = resp.State.Set(ctx, target)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *dpsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = relytClient
}
