package riskconsole

import (
	"fmt"
	"github.com/volcengine/volc-sdk-golang/base"
	"net/http"
	"net/url"
	"time"
)

type RiskConsole struct {
	*base.Client
	retry       bool
	serviceInfo *base.ServiceInfo
	apiInfo     map[string]*base.ApiInfo
}

func NewRiskConsoleInstance() *RiskConsole {
	instance := &RiskConsole{
		Client: base.NewClient(ServiceInfoMap[base.RegionCnNorth1], ApiInfoList),
		retry:  false,
	}
	return instance
}

func GetServiceInfo(region string, host string, timeout time.Duration) *base.ServiceInfo {
	return &base.ServiceInfo{
		Timeout: timeout,
		Host:    host,
		Header: http.Header{
			"Accept": []string{"application/json"},
		},
		Scheme:      "http",
		Credentials: base.Credentials{Region: region, Service: "risk_console"},
	}
}

func (p *RiskConsole) Retry() bool {
	return p.retry
}

func (p *RiskConsole) CloseRetry() {
	p.retry = false
}

func (p *RiskConsole) SetRegion(region string) error {
	serviceInfo, ok := ServiceInfoMap[region]
	if !ok {
		return fmt.Errorf("region does not spport or unknown region")
	}
	p.ServiceInfo = serviceInfo
	p.SetScheme("http")
	return nil
}

var (
	ServiceInfoMap = map[string]*base.ServiceInfo{
		base.RegionCnNorth1: {
			Timeout: 5 * time.Second,
			Host:    "riskcontrol.volcengineapi.com",
			Header: http.Header{
				"Accept": []string{"application/json"},
			},
			Scheme:      "http",
			Credentials: base.Credentials{Region: base.RegionCnNorth1, Service: "risk_console"},
		},
	}

	ApiInfoList = map[string]*base.ApiInfo{
		"GetUploadId": {
			Method:  http.MethodPost,
			Timeout: 10 * time.Second,
			Path:    "/",
			Query: url.Values{
				"Action":  []string{"GetUploadId"},
				"Version": []string{"2022-12-23"},
			},
		},
		"GetUploadedPartList": {
			Method:  http.MethodGet,
			Timeout: 10 * time.Second,
			Path:    "/",
			Query: url.Values{
				"Action":  []string{"GetUploadedPartList"},
				"Version": []string{"2022-12-23"},
			},
		},
		"UploadFile": {
			Method:  http.MethodPost,
			Timeout: 10 * time.Second,
			Path:    "/",
			Query: url.Values{
				"Action":  []string{"UploadFile"},
				"Version": []string{"2022-12-23"},
			},
		},
		"CompleteUploadFile": {
			Method:  http.MethodPost,
			Timeout: 10 * time.Second,
			Path:    "/",
			Query: url.Values{
				"Action":  []string{"CompleteUploadFile"},
				"Version": []string{"2022-12-23"},
			},
		},
	}
)
