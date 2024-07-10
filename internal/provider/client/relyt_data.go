package client

/**
{
    "alias": "zhanlu-01",
    "domain": "zhanlu",
    "variant": {
        "id": "basic"
    },
    "defaultDps": {
        "name": "hybrid",
        "description": "default dps",
        "engine": "hybrid",
        "spec": {
            "id": 2
        }
    },
    "edition": {
        "id": "standard"
    },
    "region": {
        "cloud":{
            "id": "ksc"
        },
        "id": "beijing-cicd"
    }
}
*/

const (
	DWSU_STATUS_READY = "READY"
)

type CommonRelytResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type Creator struct {
	CreateTimestamp int64  `json:"createTimestamp"`
	Domain          string `json:"domain"`
	Email           string `json:"email"`
	ID              string `json:"id"`
	IsRoot          bool   `json:"isRoot"`
	RoleID          string `json:"roleId"`
	RootAccountID   string `json:"rootAccountId"`
	Status          string `json:"status"`
}
type UsageRates struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}
type AqsSpec struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	UsageRates []UsageRates `json:"usageRates"`
}
type Owner struct {
	CreateTimestamp int64  `json:"createTimestamp"`
	Domain          string `json:"domain"`
	Email           string `json:"email"`
	ID              string `json:"id"`
	IsRoot          bool   `json:"isRoot"`
	RoleID          string `json:"roleId"`
	RootAccountID   string `json:"rootAccountId"`
	Status          string `json:"status"`
}
type Spec struct {
	ID         int64        `json:"id"`
	Name       string       `json:"name"`
	UsageRates []UsageRates `json:"usageRates"`
}
type DpsMode struct {
	AqsSpec                    AqsSpec `json:"aqsSpec"`
	CreateTime                 int64   `json:"createTime"`
	Creator                    Creator `json:"creator"`
	Description                string  `json:"description"`
	EnableAdaptiveQueryScaling bool    `json:"enableAdaptiveQueryScaling"`
	EnableAutoResume           bool    `json:"enableAutoResume"`
	EnableAutoSuspend          bool    `json:"enableAutoSuspend"`
	Engine                     string  `json:"engine"`
	ID                         string  `json:"id"`
	KeepAliveTime              int     `json:"keepAliveTime"`
	Name                       string  `json:"name"`
	Owner                      Owner   `json:"owner"`
	Spec                       Spec    `json:"spec"`
	Status                     string  `json:"status"`
	UpdateTime                 int64   `json:"updateTime"`
}

type Features struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
}
type Edition struct {
	Description string     `json:"description"`
	Features    []Features `json:"features"`
	ID          string     `json:"id"`
	IsAvailable bool       `json:"isAvailable"`
	Name        string     `json:"name"`
}
type AdditionalProp1 struct {
}
type AdditionalProp2 struct {
}
type AdditionalProp3 struct {
}
type Extensions struct {
	AdditionalProp1 AdditionalProp1 `json:"additionalProp1"`
	AdditionalProp2 AdditionalProp2 `json:"additionalProp2"`
	AdditionalProp3 AdditionalProp3 `json:"additionalProp3"`
}
type Endpoints struct {
	Extensions Extensions `json:"extensions"`
	Host       string     `json:"host"`
	ID         string     `json:"id"`
	Open       bool       `json:"open"`
	Port       int        `json:"port"`
	Protocol   string     `json:"protocol"`
	Type       string     `json:"type"`
	URI        string     `json:"uri"`
}
type Cloud struct {
	ID          string `json:"id"`
	IsAvailable bool   `json:"isAvailable"`
	IsPublic    bool   `json:"isPublic"`
	Link        string `json:"link"`
	Name        string `json:"name"`
}
type Info struct {
}
type RegionInfo struct {
	Info Info `json:"info"`
}
type Region struct {
	Area       string     `json:"area"`
	Cloud      Cloud      `json:"cloud"`
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Public     bool       `json:"public"`
	RegionInfo RegionInfo `json:"regionInfo"`
}
type Variant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type DwsuModel struct {
	Alias           string      `json:"alias"`
	CreateTimestamp int64       `json:"createTimestamp"`
	Creator         Creator     `json:"creator"`
	DefaultDps      DpsMode     `json:"defaultDps"`
	Domain          string      `json:"domain"`
	Edition         Edition     `json:"edition"`
	Endpoints       []Endpoints `json:"endpoints"`
	ID              string      `json:"id"`
	Owner           Owner       `json:"owner"`
	Region          Region      `json:"region"`
	Status          string      `json:"status"`
	Tags            []string    `json:"tags"`
	UpdateTimestamp int64       `json:"updateTimestamp"`
	Variant         Variant     `json:"variant"`
}
type CommonPage[T any] struct {
	PageNumber int  `json:"pageNumber"`
	PageSize   int  `json:"pageSize"`
	Records    []*T `json:"records"`
	Total      int  `json:"total"`
}
