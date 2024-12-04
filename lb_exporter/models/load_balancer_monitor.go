package models

import (
	"encoding/json"
	"exporter-demo/logger"
	"fmt"
)

type LoadBalancerMonitorResponse struct {
	Action     string  `json:"action"`
	ResourceID string  `json:"resource_id"`
	RetCode    int     `json:"ret_code"`
	MeterSet   []Meter `json:"meter_set"`
}

type Meter struct {
	DataSet []DataSet `json:"data_set"`
	MeterID string    `json:"meter_id"`
}

type DataSet struct {
	Data    []interface{} `json:"data"`
	EipID   string        `json:"eip_id"`
	NodeIdx int           `json:"node_idx"`
}

func ParseLoadBalancerMonitor(data []byte) *LoadBalancerMonitorResponse {
	var loadbalancerMonior LoadBalancerMonitorResponse

	err := json.Unmarshal(data, &loadbalancerMonior)
	if err != nil {
		logger.Error("Error parsing JSON:", err)
		fmt.Println("Error parsing JSON:", err)
	}

	return &loadbalancerMonior

}
