package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"net/url"
	"strconv"
	"testing"
)

var (
	client, _ = NewRelytClient(RelytClientConfig{ApiHost: host, AuthKey: auth, Role: role})
	ctx       = context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
)

const (
	host   = "http://120.92.213.241:8080"
	auth   = "9a3727e5b9c0ddabaGbll2HVLVKLLY1AyjOilAqeyPOBAb74A7VlJRAdTi0bJWJd"
	role   = "343842875420708874"
	region = "http://120.92.110.101:80"
)

func TestCreateDwsu(t *testing.T) {
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	request := DwsuModel{
		Alias:  "qingdeng-test",
		Domain: "qqqq-tst",
		Variant: &Variant{
			ID: "basic",
		},
		DefaultDps: &DpsMode{
			Name:        "hybrid",
			Description: "qingdeng-test",
			Engine:      "hybrid",
			Spec: &Spec{
				ID: 2,
			},
		},
		Edition: &Edition{
			ID: "standard",
		},
		Region: &Region{
			Cloud: &Cloud{
				ID: "ksc",
			},
			ID: "beijing-cicd",
		},
	}
	dwsu, err := client.CeateDwsu(ctx, request)
	fmt.Println(fmt.Sprintf("create result:%s resp:%s", strconv.FormatBool(err != nil), dwsu.Msg))

}

func TestListSpec(t *testing.T) {
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	spec, err := client.ListSpec(ctx, "standard", "hybrid", "ksc", "beijing-cicd")
	if err != nil {
		fmt.Println("get error" + err.Error())
	}
	marshal, err := json.Marshal(spec)
	if err != nil {
		return
	}
	fmt.Println("spec list:" + string(marshal))
}

func TestListDwsu(t *testing.T) {
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	response, err := client.ListDwsu(ctx, 100, 1)
	marshal, _ := json.Marshal(response)
	fmt.Println(fmt.Sprintf("list result:%t resp:%s", err == nil, string(marshal)))

}

func TestDeleteDwsu(t *testing.T) {
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	err := client.DropDwsu(ctx, "4679216502528")
	if err != nil {
		fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
}

func TestDeleteDps(t *testing.T) {
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	err := client.DropEdps(ctx, region, "4679367072512", "4679367072512-1472")
	if err != nil {
		fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
}

func TestPath(t *testing.T) {
	//path := client.ApiHost + "/qingdeng@zbyte-inc.com"
	path := client.ApiHost + "/中午@zbyte-inc.com"
	escape := url.PathEscape(path)
	fmt.Println(escape)
}

func TestCreateAccount(t *testing.T) {
	client.RegionApi = region
	account, err := client.CreateAccount(ctx, "", "4679350645248", Account{
		InitPassword: "zZefE#12344R*",
		Name:         "edit123",
	})
	if err != nil {
		println("create account: " + err.Error())
		return
	}
	marshal, err := json.Marshal(account)
	fmt.Println("create result: " + string(marshal))
}

func TestDropAccount(t *testing.T) {
	client.RegionApi = region
	err := client.DropAccount(ctx, region, "4679367072512", "demo3")
	if err != nil {
		println("delete account: " + err.Error())
		return
	}
}

func TestGetOpenApiMeta(t *testing.T) {
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	//client.ApiHost = "http://k8s-zbyteapp-bluewhal-254555b2ab-f9c4321b7ca3d6b0.elb.ap-east-1.amazonaws.com"
	client.ApiHost = "http://120.92.213.241:8080"
	//client.AuthKey = "801b901dce1a98f2QCH6yiakoTAVMgF0ssLc2tjvJ5duk0s5sa4j919DIBfkiCxd"
	client.AuthKey = "9a3727e5b9c0ddabaGbll2HVLVKLLY1AyjOilAqeyPOBAb74A7VlJRAdTi0bJWJd"
	//client.Role = ""
	//meta, err := client.GetOpenApiMeta(ctx, "aws", "ap-east-1")
	meta, err := client.GetOpenApiMeta(ctx, "ksc", "beijing-cicd")
	if err != nil {
		fmt.Println(fmt.Sprintf("get dwsu%s", err.Error()))
	}
	marshal, err := json.Marshal(meta)
	if err != nil {
		fmt.Println(fmt.Sprintf("err get %s", err.Error()))
		return
	}
	fmt.Println(fmt.Sprintf("get dwsu%s", string(marshal)))

}

func TestGetDwsu(t *testing.T) {
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	mode, err := client.GetDwsu(ctx, "4677306879744")
	if err != nil {
		fmt.Println(fmt.Sprintf("get dwsu%s", err.Error()))
	}
	marshal, err := json.Marshal(mode)
	if err != nil {
		fmt.Println(fmt.Sprintf("err get %s", err.Error()))
		return
	}
	fmt.Println(fmt.Sprintf("get dwsu%s", string(marshal)))
}

func TestGetDwsuApiMeta(t *testing.T) {
	mode, err := client.GetDwsuOpenApiMeta(ctx, "4677306879744")
	if err != nil {
		fmt.Println(fmt.Sprintf("get api meta%s", err.Error()))
	}
	marshal, err := json.Marshal(mode)
	if err != nil {
		fmt.Println(fmt.Sprintf("err get %s", err.Error()))
		return
	}
	fmt.Println(fmt.Sprintf("get api meta %s", string(marshal)))
}

func TestGetDps(t *testing.T) {
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	mode, err := client.GetDps(ctx, "http://120.92.110.101:80", "4679350645248", "4679350645248-1458")
	if err != nil {
		fmt.Println(fmt.Sprintf("get dps error %s", err.Error()))
	}
	marshal, err := json.Marshal(mode)
	if err != nil {
		fmt.Println(fmt.Sprintf("err get %s", err.Error()))
		return
	}
	fmt.Println(fmt.Sprintf("get dps %s", string(marshal)))
}
