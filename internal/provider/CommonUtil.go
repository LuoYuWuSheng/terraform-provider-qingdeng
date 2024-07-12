package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"terraform-provider-relyt/internal/provider/client"
)

func RouteRegionUri(ctx context.Context, dwsuId string, client *client.RelytClient,
	diag *diag.Diagnostics) *client.OpenApiMetaInfo {
	meta, err := client.GetDwsuOpenApiMeta(ctx, dwsuId)
	if err != nil || meta == nil {
		errMsg := "get RegionApi is nil"
		if err != nil {
			errMsg = err.Error()
		}
		diag.AddError("error get region api", "fail to get Region uri address dwsuID:"+
			""+dwsuId+" error: "+errMsg)
		return meta
	}
	return meta
}
