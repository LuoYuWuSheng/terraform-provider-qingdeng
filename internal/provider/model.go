package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DefaultDps struct {
	Description types.String `tfsdk:"description"`
	Engine      types.String `tfsdk:"engine"`
	Name        types.String `tfsdk:"name"`
	Size        types.Int64  `tfsdk:"size"`
	Status      types.String `tfsdk:"status"`
}

type DwsuModel struct {
	ID          types.String `tfsdk:"id"`
	Alias       types.String `tfsdk:"alias"`
	Cloud       types.String `tfsdk:"cloud"`
	Region      types.String `tfsdk:"region"`
	DwsuType    types.String `tfsdk:"dwsu_type"`
	DefaultDps  DefaultDps   `tfsdk:"defaultDps"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

type DpsModel struct {
	Dwsu DwsuModel `tfsdk:"dwsu"`

	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Engine      types.String `tfsdk:"engine"`
	Size        types.Int64  `tfsdk:"size"`
	LastUpdated types.String `tfsdk:"last_updated"`
}
