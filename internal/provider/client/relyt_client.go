package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func NewRelytClient(apiHost, authKey, roleId string) RelytClient {
	return RelytClient{apiHost: apiHost, authKey: authKey, roleId: roleId}
}

type RelytClient struct {
	apiHost string
	authKey string
	roleId  string
}

func (p *RelytClient) ListDwsu(ctx context.Context, pageSize, pageNumber int) ([]*DwsuModel, error) {
	resp := CommonRelytResponse[CommonPage[DwsuModel]]{}
	pageQuery := map[string]string{
		"pageSize":   strconv.Itoa(pageSize),
		"pageNumber": strconv.Itoa(pageNumber),
	}
	err := doHttpRequest(p, ctx, "/dwsu", "GET", &resp, nil, pageQuery)
	if err != nil {
		return nil, err
	}
	return resp.Data.Records, nil
}

func (p *RelytClient) CeateDwsu(ctx context.Context, request DwsuModel) (*CommonRelytResponse[string], error) {
	url := "/dwsu"
	resp := CommonRelytResponse[string]{}
	err := doHttpRequest(p, ctx, url, "POST", &resp, request, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *RelytClient) GetDwsu(ctx context.Context, dwServiceUnitId string) (*DwsuModel, error) {
	path := fmt.Sprintf("/dwsu/%s", dwServiceUnitId)
	resp := CommonRelytResponse[DwsuModel]{}
	err := doHttpRequest(p, ctx, path, "GET", &resp, nil, nil)
	if err != nil {
		tflog.Error(ctx, "Error get dwsu:"+err.Error())
		return nil, err
	}
	return &resp.Data, nil
}

func (p *RelytClient) DropDwsu(ctx context.Context, dwServiceUnitId string) error {
	path := fmt.Sprintf("/dwsu/%s", dwServiceUnitId)
	resp := CommonRelytResponse[string]{}
	err := doHttpRequest(p, ctx, path, "DELETE", &resp, nil, nil)
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
	err := doHttpRequest(p, ctx, path, "GET", &resp, nil, pageQuery)
	if err != nil {
		return nil, err
	}
	return resp.Data.Records, nil
}

func (p *RelytClient) CreateEdps(ctx context.Context, dwServiceUnitId string, mode DpsMode) (*CommonRelytResponse[string], error) {
	path := fmt.Sprintf("/dwsu/%s/dps", dwServiceUnitId)
	resp := CommonRelytResponse[string]{}
	if err := doHttpRequest(p, ctx, path, "POST", &resp, mode, nil); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *RelytClient) GetDps(ctx context.Context, dwServiceUnitId, dpsBizId string) (*DpsMode, error) {
	path := fmt.Sprintf("/dwsu/%s/dps/%s", dwServiceUnitId, dpsBizId)
	resp := CommonRelytResponse[DpsMode]{}
	err := doHttpRequest(p, ctx, path, "GET", &resp, nil, nil)
	if err != nil {
		tflog.Error(ctx, "Error get dps:"+err.Error())
		return nil, err
	}
	return &resp.Data, nil
}

func (p *RelytClient) DropEdps(ctx context.Context, dwServiceUnitId, dpsBizId string) error {
	path := fmt.Sprintf("/dwsu/%s/dps/%s", dwServiceUnitId, dpsBizId)
	resp := CommonRelytResponse[string]{}
	err := doHttpRequest(p, ctx, path, "DELETE", &resp, nil, nil)
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
	err := doHttpRequest(p, ctx, path, "GET", &specList, nil, parameter)
	if err != nil {
		return nil, err
	}
	return specList.Data, nil
}

func (p *RelytClient) CreateAccount(ctx context.Context) {
	path := "/dwsu/{dwServiceUnitId}/account"
	print(path)
}

func (p *RelytClient) AsyncAccountConfig() {
	path := "/dwsu/{dwServiceUnitId}/user/{userId}/asyncresult"
	print(path)
}

func doHttpRequest[T any](p *RelytClient, ctx context.Context, path, method string, respMode *CommonRelytResponse[T], request any, parameter map[string]string) (err error) {

	var jsonData = []byte("")
	if request != nil && "" != request {
		requestJson, err := json.Marshal(request)
		if err != nil {
			tflog.Error(ctx, "fmt request json error:"+err.Error())
		}
		tflog.Info(ctx, "request data :"+string(requestJson))
		jsonData = requestJson // POST请求发送的数据
	}

	hostApi := p.apiHost + path
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
	req.Header.Set("x-maxone-api-key", p.authKey)
	req.Header.Set("x-maxone-role-id", p.roleId)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tflog.Error(ctx, "Error sending request:"+err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		tflog.Error(ctx, "Error status http code not 200! "+resp.Status)
		return fmt.Errorf("Error status http code not 200! " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tflog.Error(ctx, "Error reading response body:"+err.Error())
		return err
	}
	tflog.Info(ctx, "Response:"+string(body))

	err = json.Unmarshal(body, respMode)
	if err != nil {
		tflog.Error(ctx, "read json respFail:"+err.Error())
		return err
	}
	if respMode.Code != 200 {
		tflog.Error(ctx, "error call api! resp code not 200: "+string(body))
		return fmt.Errorf(string(body))
	}
	return nil
}

func (p *RelytClient) GetDwsuByAlias(ctx context.Context, alias string) (*DwsuModel, error) {
	models, err := p.ListDwsu(ctx, 100, 1)
	if err == nil && len(models) > 0 {
		for _, model := range models {
			if model.Alias == alias {
				return model, nil
			}
		}
	}
	return nil, fmt.Errorf("can't find id")
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

func printResp(ctx context.Context, resp http.Response) string {
	// 打印响应状态码
	var builder strings.Builder
	m := "Status Code:" + strconv.Itoa(resp.StatusCode)
	tflog.Info(ctx, m)
	builder.WriteString(m + "\n")

	// 打印响应头信息
	tflog.Info(ctx, "Headers:")
	for key, value := range resp.Header {
		val := key + " :" + strings.Join(value, ",")
		tflog.Info(ctx, val)
		builder.WriteString(val + "\n")
	}

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tflog.Info(ctx, "Error reading response body:"+err.Error())
		return err.Error()
	}
	tflog.Info(ctx, "Body:"+string(body))
	builder.WriteString(string(body))
	return builder.String()
}
