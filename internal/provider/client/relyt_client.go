package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

func NewRelytClient(config RelytClientConfig) (RelytClient, error) {
	return RelytClient{config}, nil
}

type RelytClientConfig struct {
	ApiHost   string
	AuthKey   string
	Role      string
	RegionApi string
}

type RelytClient struct {
	RelytClientConfig
}

func (p *RelytClient) ListDwsu(ctx context.Context, pageSize, pageNumber int) ([]*DwsuModel, error) {
	resp := CommonRelytResponse[CommonPage[DwsuModel]]{}
	pageQuery := map[string]string{
		"pageSize":   strconv.Itoa(pageSize),
		"pageNumber": strconv.Itoa(pageNumber),
	}
	err := doHttpRequest(p, ctx, "", "/dwsu", "GET", &resp, nil, pageQuery, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data.Records, nil
}

func (p *RelytClient) CeateDwsu(ctx context.Context, request DwsuModel) (*CommonRelytResponse[string], error) {
	url := "/dwsu"
	resp := CommonRelytResponse[string]{}
	err := doHttpRequest(p, ctx, "", url, "POST", &resp, request, nil, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *RelytClient) GetDwsu(ctx context.Context, dwServiceUnitId string) (*DwsuModel, error) {
	path := fmt.Sprintf("/dwsu/%s", dwServiceUnitId)
	resp := CommonRelytResponse[DwsuModel]{}
	err := doHttpRequest(p, ctx, "", path, "GET", &resp, nil, nil, nil)
	if err != nil {
		tflog.Error(ctx, "Error get dwsu:"+err.Error())
		return nil, err
	}
	return resp.Data, nil
}

func (p *RelytClient) DropDwsu(ctx context.Context, dwServiceUnitId string) error {
	path := fmt.Sprintf("/dwsu/%s", dwServiceUnitId)
	resp := CommonRelytResponse[string]{}
	err := doHttpRequest(p, ctx, "", path, "DELETE", &resp, nil, nil, nil)
	if err != nil {
		tflog.Info(ctx, "delete dwsu err:"+err.Error())
		return err
	}
	return nil
}

func (p *RelytClient) ListDps(ctx context.Context, pageSize, pageNumber int, dwServiceUnitId string) ([]*DpsMode, error) {
	resp := CommonRelytResponse[CommonPage[DpsMode]]{}
	pageQuery := map[string]string{
		"pageSize":   strconv.Itoa(pageSize),
		"pageNumber": strconv.Itoa(pageNumber),
	}
	path := fmt.Sprintf("/dwsu/%s/dps", dwServiceUnitId)
	err := doHttpRequest(p, ctx, "", path, "GET", &resp, nil, pageQuery, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data.Records, nil
}

func (p *RelytClient) CreateEdps(ctx context.Context, regionUri string, dwServiceUnitId string, mode DpsMode) (*CommonRelytResponse[string], error) {
	path := fmt.Sprintf("/dwsu/%s/dps", dwServiceUnitId)
	resp := CommonRelytResponse[string]{}
	if err := doHttpRequest(p, ctx, regionUri, path, "POST", &resp, mode, nil, nil); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *RelytClient) GetDps(ctx context.Context, regionUri, dwServiceUnitId, dpsBizId string) (*DpsMode, error) {
	path := fmt.Sprintf("/dwsu/%s/dps/%s", dwServiceUnitId, dpsBizId)
	resp := CommonRelytResponse[DpsMode]{}
	err := doHttpRequest(p, ctx, regionUri, path, "GET", &resp, nil, nil, nil)
	if err != nil {
		tflog.Error(ctx, "Error get dps:"+err.Error())
		return nil, err
	}
	return resp.Data, nil
}

func (p *RelytClient) DropEdps(ctx context.Context, regionUri, dwServiceUnitId, dpsBizId string) error {
	path := fmt.Sprintf("/dwsu/%s/dps/%s", dwServiceUnitId, dpsBizId)
	resp := CommonRelytResponse[string]{}
	handler := func(response *CommonRelytResponse[string], respString []byte) (*CommonRelytResponse[string], error) {
		if response.Code != 200 && resp.Code != CODE_DPS_NOT_FOUND {
			body := string(respString)
			tflog.Error(ctx, "error call api! resp code not success! body: "+body)
			return response, fmt.Errorf(body)
		}
		return nil, nil
	}
	err := doHttpRequest(p, ctx, regionUri, path, "DELETE", &resp, nil, nil, handler)
	if err != nil {
		tflog.Info(ctx, "delete dps err:"+err.Error())
		return err
	}
	return nil
}

func (p *RelytClient) ListSpec(ctx context.Context, edition, dpsType, cloud, region string) ([]Spec, error) {
	path := fmt.Sprintf("/dwsu/edition/%s/dps/%s/specs", edition, dpsType)
	specList := CommonRelytResponse[[]Spec]{}
	parameter := map[string]string{"cloud": cloud, "region": region}
	err := doHttpRequest(p, ctx, "", path, "GET", &specList, nil, parameter, nil)
	if err != nil {
		return nil, err
	}
	return *specList.Data, nil
}

func (p *RelytClient) CreateAccount(ctx context.Context, regionUri string, dwsuId string, account Account) (*CommonRelytResponse[string], error) {
	path := fmt.Sprintf("/dwsu/%s/account", dwsuId)
	resp := CommonRelytResponse[string]{}
	err := doHttpRequest(p, ctx, regionUri, path, "POST", &resp, account, nil, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *RelytClient) DropAccount(ctx context.Context, regionUri string, dwsuId string, userId string) error {
	path := fmt.Sprintf("/dwsu/%s/user/%s", dwsuId, url.PathEscape(userId))
	resp := CommonRelytResponse[string]{}
	handler := func(response *CommonRelytResponse[string], respString []byte) (*CommonRelytResponse[string], error) {
		if response.Code != 200 && resp.Code != CODE_USER_NOT_FOUND {
			body := string(respString)
			tflog.Error(ctx, "error call api! resp code not success! body: "+body)
			return response, fmt.Errorf(body)
		}
		return nil, nil
	}
	err := doHttpRequest(p, ctx, regionUri, path, "DELETE", &resp, nil, nil,
		handler)
	return err
}

func (p *RelytClient) AsyncAccountConfig(ctx context.Context, regionUri, dwsuId, userId string, asyncResult AsyncResult) (*CommonRelytResponse[string], error) {
	path := fmt.Sprintf("/dwsu/%s/user/%s/asyncresult", dwsuId, url.PathEscape(userId))
	resp := CommonRelytResponse[string]{}
	err := doHttpRequest(p, ctx, regionUri, path, "PUT", &resp, asyncResult, nil, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *RelytClient) LakeFormationConfig(ctx context.Context, regionUri, dwsuId, userId string, formation LakeFormation) (*CommonRelytResponse[string], error) {
	path := fmt.Sprintf("/dwsu/%s/user/%s/lakeformation", dwsuId, url.PathEscape(userId))
	resp := CommonRelytResponse[string]{}
	err := doHttpRequest(p, ctx, regionUri, path, "PUT", &resp, formation, nil, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *RelytClient) GetOpenApiMeta(ctx context.Context, cloud, region string) (*OpenApiMetaInfo, error) {
	//fixme uri编码需要，否则发送请求不会支持
	path := fmt.Sprintf("/infra/%s/%s/endpoint", url.PathEscape(cloud), url.PathEscape(region))
	resp := CommonRelytResponse[[]*OpenApiMetaInfo]{}
	err := doHttpRequest(p, ctx, "", path, "GET", &resp, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	lengthOfApi := len(*resp.Data)
	if lengthOfApi != 1 {
		return nil, fmt.Errorf("error read regionApi! length of api " + strconv.Itoa(lengthOfApi))
	}
	return (*resp.Data)[0], nil
}

func (p *RelytClient) GetDwsuOpenApiMeta(ctx context.Context, dwsuId string) (*OpenApiMetaInfo, error) {
	dwsu, err := p.GetDwsu(ctx, dwsuId)
	if err != nil {
		return nil, err
	}
	if dwsu == nil {
		return nil, fmt.Errorf("can't find dwsu meta! %s", dwsuId)
	}
	meta, err := p.GetOpenApiMeta(ctx, dwsu.Region.Cloud.ID, dwsu.Region.ID)
	return meta, err
}

func doHttpRequest[T any](p *RelytClient, ctx context.Context, host, path, method string,
	respMode *CommonRelytResponse[T],
	request any,
	parameter map[string]string,
	codeHandler func(response *CommonRelytResponse[T], respDumpByte []byte) (*CommonRelytResponse[T], error)) (err error) {
	if host == "" {
		host = p.ApiHost
	}
	var jsonData = []byte("")
	if request != nil && "" != request {
		requestJson, err := json.Marshal(request)
		if err != nil {
			tflog.Error(ctx, "fmt request json error:"+err.Error())
		}
		tflog.Info(ctx, "request data :"+string(requestJson))
		jsonData = requestJson // POST请求发送的数据
	}
	hostApi := host + path
	parsedHostApi, err := url.Parse(hostApi)
	if err != nil {
		return err
	}
	queryParams := url.Values{}
	if parameter != nil {
		for k, v := range parameter {
			queryParams.Add(k, v)
		}
	}
	parsedHostApi.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(method, parsedHostApi.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		tflog.Error(ctx, "Error creating request:"+err.Error())
		return err
	}
	req.Header.Set("x-maxone-api-key", p.AuthKey)
	req.Header.Set("x-maxone-role-id", p.Role)
	req.Header.Set("Content-Type", "application/json")
	requestString, _ := httputil.DumpRequestOut(req, true)
	tflog.Info(ctx, "== request: "+string(requestString))
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		tflog.Error(ctx, "Error sending request:"+err.Error())
		return err
	}
	responseString, _ := httputil.DumpResponse(resp, true)
	tflog.Info(ctx, "== response: "+string(responseString))
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		tflog.Error(ctx, "Error status http code not 200! "+resp.Status)
		//printResp(ctx, resp)
		return fmt.Errorf("Error status http code not 200! %s resp: %s ", resp.Status, string(responseString))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tflog.Error(ctx, "Error reading responseString body:"+err.Error())
		return err
	}
	//tflog.Info(ctx, "Response:"+string(body))

	err = json.Unmarshal(body, respMode)
	if err != nil {
		tflog.Error(ctx, "read json respFail:"+err.Error())
		return err
	}
	if respMode.Code != 200 {
		tflog.Warn(ctx, "error call api! resp code not 200: "+string(body))
	}
	if codeHandler != nil {
		tflog.Trace(ctx, "use code handle func!")
		handler, err := codeHandler(respMode, responseString)
		if handler != nil {
			respMode.Code = handler.Code
			respMode.Data = handler.Data
			respMode.Msg = handler.Msg
		}
		if err != nil {
			return err
		}
	} else {
		if respMode.Code != 200 {
			tflog.Error(ctx, "error call api! resp code not 200: "+string(body))
			return fmt.Errorf(string(body))
		}
	}
	return nil
}

func (p *RelytClient) TimeOutTask(timeoutSec int, task func() (any, error)) (any, error) {
	// 设置超时时间
	timeout := time.Duration(timeoutSec) * time.Second

	// 创建带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	//// 启动任务
	//done := make(chan bool)
	//f := func() (any, error) {
	//
	//}
	//go f
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task timed out")
			//done <- false
			return nil, fmt.Errorf("timeout")
		default:
			a, err := task()
			if err == nil {
				//done <- true
				return a, err
			}
			time.Sleep(5 * time.Second)
		}
	}
}
