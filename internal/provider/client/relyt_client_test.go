package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"strconv"
	"testing"
)

var client = NewRelytClient(host, auth, role)

const (
	host = "http://120.92.213.241:8080"
	auth = "9a3727e5b9c0ddabaGbll2HVLVKLLY1AyjOilAqeyPOBAb74A7VlJRAdTi0bJWJd"
	role = "343842875420708874"
)

func TestCreateDwsu(t *testing.T) {
	client := NewRelytClient(host, auth, role)
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	request := DwsuModel{
		Alias:  "qingdeng-test",
		Domain: "qqqq-tst",
		Variant: Variant{
			ID: "basic",
		},
		DefaultDps: DpsMode{
			Name:        "hybrid",
			Description: "qingdeng-test",
			Engine:      "hybrid",
			Spec: Spec{
				ID: 2,
			},
		},
		Edition: Edition{
			ID: "standard",
		},
		Region: Region{
			Cloud: Cloud{
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
	client := NewRelytClient(host, auth, role)
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	response, err := client.ListDwsu(ctx, 100, 1)
	marshal, _ := json.Marshal(response)
	fmt.Println(fmt.Sprintf("list result:%t resp:%s", err == nil, string(marshal)))

}

func TestDeleteDwsu(t *testing.T) {
	client := NewRelytClient(host, auth, role)
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	err := client.DropDwsu(ctx, "4679216502528")
	if err != nil {
		fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
}

func TestTimeOutRead(t *testing.T) {
	client := NewRelytClient(host, auth, role)
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	resp, err := client.TimeOutTask(3, func() (any, error) {
		return client.GetDwsuByAlias(ctx, "bar")
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
	if resp != nil {
		model := resp.(*DwsuModel)
		marshal, err := json.Marshal(model)
		if err != nil {
			fmt.Println(fmt.Sprintf("json %s", err.Error()))
		}
		fmt.Println(string(marshal))
	}
}

func TestGetDwsu(t *testing.T) {
	client := NewRelytClient(host, auth, role)
	ctx := context.WithValue(context.Background(), "provider", hclog.NewInterceptLogger)
	mode, err := client.GetDwsu(ctx, "4679216502528")
	if err != nil {
		fmt.Println(fmt.Sprintf("drop dwsu%s", err.Error()))
	}
	marshal, err := json.Marshal(mode)
	if err != nil {
		fmt.Println(fmt.Sprintf("err get %s", err.Error()))
		return
	}
	fmt.Println(fmt.Sprintf("get dwsu%s", string(marshal)))
}
