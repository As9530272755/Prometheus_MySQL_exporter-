package models

import (
	"encoding/json"
	"exporter-demo/logger"
	"fmt"
)

type LoadbalancerSet struct {
	NodeBalanceMode  string `json:"node_balance_mode"`
	BackendIPVersion int    `json:"backend_ip_version"`
	IsApplied        int    `json:"is_applied"`
	VxnetID          string `json:"vxnet_id"`
	ConsoleID        string `json:"console_id"`
	Cluster          []struct {
		EipName   string `json:"eip_name,omitempty"`
		EipAddr   string `json:"eip_addr,omitempty"`
		EipID     string `json:"eip_id,omitempty"`
		Instances []struct {
			InstanceID string `json:"instance_id"`
			VgwMgmtIP  string `json:"vgw_mgmt_ip"`
			NodeIdx    int    `json:"node_idx"`
		} `json:"instances"`
	} `json:"cluster"`
	CreateTime     string `json:"create_time"`
	Rsyslog        string `json:"rsyslog"`
	Owner          string `json:"owner"`
	PlaceGroupID   string `json:"place_group_id"`
	SecurityGroups []struct {
		GroupID   string `json:"group_id"`
		GroupName string `json:"group_name"`
	} `json:"security_groups"`
	Features         int         `json:"features"`
	SubCode          int         `json:"sub_code"`
	RebuildLastTime  interface{} `json:"rebuild_last_time"`
	SecurityGroupID  string      `json:"security_group_id"`
	LoadbalancerType int         `json:"loadbalancer_type"`
	LoadbalancerName string      `json:"loadbalancer_name"`
	ActiveStandby    int         `json:"active_standby"`
	Memory           int         `json:"memory"`
	StatusTime       string      `json:"status_time"`
	NodeCount        int         `json:"node_count"`
	Listeners        []struct {
		Forwardfor               int           `json:"forwardfor"`
		ListenerOption           int           `json:"listener_option"`
		BalanceMode              string        `json:"balance_mode"`
		BackendProtocol          string        `json:"backend_protocol"`
		HealthyCheckMethod       string        `json:"healthy_check_method"`
		Scene                    int           `json:"scene"`
		ConsoleID                string        `json:"console_id"`
		Disabled                 int           `json:"disabled"`
		CreateTime               string        `json:"create_time"`
		WafDomainPolicies        []interface{} `json:"waf_domain_policies"`
		Owner                    string        `json:"owner"`
		HealthyCheckOption       string        `json:"healthy_check_option"`
		ListenerPortEnd          int           `json:"listener_port_end"`
		SessionSticky            string        `json:"session_sticky"`
		TLSVersion               string        `json:"tls_version"`
		ListenerProtocol         string        `json:"listener_protocol"`
		ServerCertificateID      []interface{} `json:"server_certificate_id"`
		LoadbalancerListenerName string        `json:"loadbalancer_listener_name"`
		HttpRedirectionCode      interface{}   `json:"http_redirection_code"`
		Controller               string        `json:"controller"`
		TunnelTimeout            int           `json:"tunnel_timeout"`
		LoadbalancerListenerID   string        `json:"loadbalancer_listener_id"`
		ListenerPort             int           `json:"listener_port"`
		RootUserID               string        `json:"root_user_id"`
		Timeout                  int           `json:"timeout"`
		LoadbalancerID           string        `json:"loadbalancer_id"`
	} `json:"listeners"`
	Vxnet struct {
		VxnetType  int    `json:"vxnet_type"`
		VxnetID    string `json:"vxnet_id"`
		InstanceID string `json:"instanc_id"`
		VxnetName  string `json:"vxnet_name"`
		PrivateIP  string `json:"private_ip"`
		NicID      string `json:"nic_id"`
	} `json:"vxnet"`
	Status              string        `json:"status"`
	EcmPDisabled        int           `json:"ecmp_disabled"`
	Description         interface{}   `json:"description"`
	Tags                []interface{} `json:"tags"`
	TransitionStatus    string        `json:"transition_status"`
	PrivateIPs          []string      `json:"private_ips"`
	Eips                []interface{} `json:"eips"`
	Controller          string        `json:"controller"`
	Repl                string        `json:"repl"`
	WafPg               string        `json:"waf_pg"`
	ServiceIpEnabled    int           `json:"service_ip_enabled"`
	ClusterMode         int           `json:"cluster_mode"`
	SriovNicType        int           `json:"sriov_nic_type"`
	ZoneID              string        `json:"zone_id"`
	Hypervisor          string        `json:"hypervisor"`
	CPU                 int           `json:"cpu"`
	RootUserID          string        `json:"root_user_id"`
	HttpHeaderSize      interface{}   `json:"http_header_size"`
	Mode                int           `json:"mode"`
	SSLDefaultDhParam   interface{}   `json:"ssl_default_dh_param"`
	ResourceProjectInfo []interface{} `json:"resource_project_info"`
	LoadbalancerID      string        `json:"loadbalancer_id"`
	VxnetCluster        struct {
		ClusterStatus string `json:"cluster_status"`
	} `json:"vxnet_cluster"`
}

// 定义结构体来匹配 JSON 数据的结构
type LoadBalancerResponse struct {
	Action           string            `json:"action"`
	TotalCount       int               `json:"total_count"`
	LoadbalancerSets []LoadbalancerSet `json:"loadbalancer_set"`
	RetCode          int               `json:"ret_code"`
}

func ParseLoadBalancers(data []byte) *LoadBalancerResponse {
	var loadBalancersResponse LoadBalancerResponse

	err := json.Unmarshal(data, &loadBalancersResponse)
	if err != nil {
		logger.Error("Error parsing JSON: ", err)
		fmt.Println("ParseLoadBalancers Error parsing JSON: ", err)
	}
	return &loadBalancersResponse
}
