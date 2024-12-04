package servser

import (
	bossclient "exporter-demo/client"
	"exporter-demo/exporter_config"
	"exporter-demo/logger"
	"exporter-demo/models"
	"fmt"
	"time"
)

func LoadBalancerMonitor() (*models.LoadBalancerMonitorResponse, error) {

	host, protocol, access_key, secret_key, zone, port, ignore_verify := exporter_config.ConfigBoss()

	// init client
	client, err := bossclient.NewClient(host, port, protocol, access_key, secret_key, ignore_verify)
	if err != nil {
		logger.Error("链接至青云client error", err.Error())
		fmt.Println("链接至青云client error", err.Error())
	}

	//获取当前时间并计算start_time和end_time
	now := time.Now().UTC()
	EndTime := now.Format(time.RFC3339)
	//起始时间设置为 10s 前
	StartTime := now.Add(-1 * time.Hour).Format(time.RFC3339) // Example: one hour ago

	// please modify following params to send a test request
	//zone := "bjdev"
	step := "10s"
	meters := []string{"request"}
	liprotocol := "tcp"
	subtype := "tcp"

	Listeners := Wjdev_Listeners().LoadbalancerListenerSet

	for _, Listener := range Listeners {
		resp, err := client.GetLoadBalancerMonitor(zone, Listener.ListenerID, step, liprotocol, subtype, StartTime, EndTime, meters)
		if err != nil {
			logger.Error("青云 LoadBalancers 获取 error", err.Error())
			fmt.Println("青云 LoadBalancers 获取 error", err.Error())
			//return nil, err
		}
		return models.ParseLoadBalancerMonitor([]byte(resp)), nil
	}

	// 通过单个监听器获取对应监控数据
	return nil, err
}
