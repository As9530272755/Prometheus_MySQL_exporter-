package models

import (
	"encoding/json"
	"exporter-demo/logger"
	"fmt"
)

// 定义结构体来匹配 JSON 数据的结构
type Backend struct {
	Status                  string `json:"status"`
	LoadbalancerBackendID   string `json:"loadbalancer_backend_id"`
	PortEnd                 int    `json:"port_end"`
	IsTop                   int    `json:"is_top"`
	Weight                  int    `json:"weight"`
	ResourceID              string `json:"resource_id"`
	Backup                  int    `json:"backup"`
	LoadbalancerBackendName string `json:"loadbalancer_backend_name"`
	ConsoleID               string `json:"console_id"`
	RootUserID              string `json:"root_user_id"`
	LoadbalancerPolicyID    string `json:"loadbalancer_policy_id"`
	Disabled                int    `json:"disabled"`
	Controller              string `json:"controller"`
	CreateTime              string `json:"create_time"`
	Port                    int    `json:"port"`
	Owner                   string `json:"owner"`
	LoadbalancerListenerID  string `json:"loadbalancer_listener_id"`
	NicID                   string `json:"nic_id"`
	LoadbalancerID          string `json:"loadbalancer_id"`
}

type Listener struct {
	ForwardFor          int         `json:"forwardfor"`
	ListenerOption      int         `json:"listener_option"`
	BalanceMode         string      `json:"balance_mode"`
	BackendProtocol     string      `json:"backend_protocol"`
	HealthyCheckMethod  string      `json:"healthy_check_method"`
	Scene               int         `json:"scene"`
	ConsoleID           string      `json:"console_id"`
	Disabled            int         `json:"disabled"`
	CreateTime          string      `json:"create_time"`
	WafDomainPolicies   []string    `json:"waf_domain_policies"`
	Owner               string      `json:"owner"`
	HealthyCheckOption  string      `json:"healthy_check_option"`
	ListenerPortEnd     int         `json:"listener_port_end"`
	SessionSticky       string      `json:"session_sticky"`
	Backends            []Backend   `json:"backends"`
	TLSVersion          string      `json:"tls_version"`
	ListenerProtocol    string      `json:"listener_protocol"`
	ServerCertificateID []string    `json:"server_certificate_id"`
	ListenerName        string      `json:"loadbalancer_listener_name"`
	HttpRedirectionCode interface{} `json:"http_redirection_code"`
	Controller          string      `json:"controller"`
	TunnelTimeout       int         `json:"tunnel_timeout"`
	ListenerID          string      `json:"loadbalancer_listener_id"`
	ListenerPort        int         `json:"listener_port"`
	ZoneID              string      `json:"zone_id"`
	RootUserID          string      `json:"root_user_id"`
	Timeout             int         `json:"timeout"`
	LoadbalancerID      string      `json:"loadbalancer_id"`
}

type ListenerSetResponse struct {
	Action                  string     `json:"action"`
	TotalCount              int        `json:"total_count"`
	LoadbalancerListenerSet []Listener `json:"loadbalancer_listener_set"`
	RetCode                 int        `json:"ret_code"`
}

func ParseListeners(data []byte) *ListenerSetResponse {
	var listenerSetResponse ListenerSetResponse

	err := json.Unmarshal([]byte(data), &listenerSetResponse)
	if err != nil {
		logger.Error("Error parsing JSON: ", err)
		fmt.Println("Error parsing JSON: ", err)
	}

	return &listenerSetResponse

}
