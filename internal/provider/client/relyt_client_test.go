package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"net/url"
	"os"
	"strconv"
	"testing"
)

func init() {
	file, err := os.ReadFile("../../../env/testConf.json")
	config := RelytClientConfig{}
	if err != nil {
		fmt.Println("error read test file")
		return
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("error read test file")
		return
	}
	client, _ = NewRelytClient(config)
	ctx = context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
}

var (
	client RelytClient
	ctx    context.Context
)

func TestReadConf(t *testing.T) {
	fmt.Println(client.RegionApi)
}

func TestCreateDwsu(t *testing.T) {

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
	response, err := client.ListDwsu(ctx, 100, 1)
	marshal, _ := json.Marshal(response)
	fmt.Println(fmt.Sprintf("list result:%t resp:%s", err == nil, string(marshal)))

}

func TestDeleteDwsu(t *testing.T) {
	err := client.DropDwsu(ctx, "4679216502528-abc-not-exist")
	if err != nil {
		fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
}

func TestDeleteDps(t *testing.T) {
	err := client.DropEdps(ctx, client.RegionApi, "4679367072512", "4679367072512-1472-abc")
	if err != nil {
		fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
}

func TestPath(t *testing.T) {
	//path := client.ApiHost + "/qingdeng@zbyte-inc.com"
	path := client.ApiHost + "/中午@zbyte-inc.com"
	escape := url.PathEscape(path)
	fmt.Println(escape)

	var st []byte
	sprintf := fmt.Sprintf("abc:%s %t ", string(st), st == nil)
	fmt.Println(sprintf)
}

func TestCreateAccount(t *testing.T) {
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
	err := client.DropAccount(ctx, client.RegionApi, "4679367072512", "demo3")
	if err != nil {
		println("delete account: " + err.Error())
		return
	}
}

func TestGetOpenApiMeta(t *testing.T) {

	//client.AuthKey = "801b901dce1a98f2QCH6yiakoTAVMgF0ssLc2tjvJ5duk0s5sa4j919DIBfkiCxd"
	//client.AuthKey = "9a3727e5b9c0ddabaGbll2HVLVKLLY1AyjOilAqeyPOBAb74A7VlJRAdTi0bJWJd"
	//client.Role = ""
	meta, err := client.GetOpenApiMeta(ctx, "aws", "ap-east-1")
	//meta, err := client.GetOpenApiMeta(ctx, "ksc", "beijing-cicd")
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

	mode, err := client.GetDwsu(ctx, "4679438371328")
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

func TestGetDwsuServiceAccount(t *testing.T) {
	api, err := client.GetDwsuOpenApiMeta(ctx, "4954339339520")
	mode, err := client.GetDwsuServiceAccount(ctx, api.URI, "4954339339520")
	//mode, err := client.GetDwsuServiceAccount(ctx, "4679344580864")
	if len(mode) > 0 {
		fmt.Println("accoiunt not zero")
	}
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
