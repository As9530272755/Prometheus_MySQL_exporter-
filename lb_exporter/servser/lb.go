package servser

import (
	bossclient "exporter-demo/client"
	"exporter-demo/exporter_config"
	"exporter-demo/logger"
	"exporter-demo/models"
	"fmt"
)

// 由于青云API每次只能获取分页 100 个数据单元，所以这里拿的是第一页前一百个LB的数据
func Wjdev_lb_set1() (*models.LoadBalancerResponse, error) {

	host, protocol, access_key, secret_key, zone, port, ignore_verify := exporter_config.ConfigBoss()

	// init client
	client, err := bossclient.NewClient(host, port, protocol, access_key, secret_key, ignore_verify)
	if err != nil {
		logger.Error("链接至青云client error", err.Error())
		fmt.Println("链接至青云client error", err.Error())
		return nil, err
	}

	// please modify following params to send a test request
	ZONE := zone
	STATUS := []string{"active", "stopped"}
	LIMIT := 99
	OFFSET := 0

	resp, err := client.Boss2DescribeLoadBalancers(ZONE, STATUS, LIMIT, OFFSET)
	if err != nil {
		logger.Error("青云 LoadBalancers 获取 error", err.Error())
		fmt.Println("链接至青云 LoadBalancers error", err.Error())
		return nil, err
	}
	return models.ParseLoadBalancers([]byte(resp)), nil

}

func Wjdev_lb_set2() (*models.LoadBalancerResponse, error) {

	host, protocol, access_key, secret_key, zone, port, ignore_verify := exporter_config.ConfigBoss()

	// init client
	client, err := bossclient.NewClient(host, port, protocol, access_key, secret_key, ignore_verify)
	if err != nil {
		logger.Error("链接至青云client error", err.Error())
		fmt.Println("链接至青云 client error", err.Error())
		return nil, err
	}

	// please modify following params to send a test request
	ZONE := zone
	STATUS := []string{"active", "stopped"}
	LIMIT := 199
	OFFSET := 100

	resp, err := client.Boss2DescribeLoadBalancers(ZONE, STATUS, LIMIT, OFFSET)
	if err != nil {
		if err != nil {
			logger.Error("青云 LoadBalancers 获取 error", err.Error())
			fmt.Println("链接至青云 LoadBalancers 获取 error", err.Error())
			return nil, err
		}
	}
	return models.ParseLoadBalancers([]byte(resp)), nil
}

func Wjdev_lb_set3() (*models.LoadBalancerResponse, error) {
	host, protocol, access_key, secret_key, zone, port, ignore_verify := exporter_config.ConfigBoss()

	// init client
	client, err := bossclient.NewClient(host, port, protocol, access_key, secret_key, ignore_verify)
	if err != nil {
		logger.Error("链接至青云client error", err.Error())
		fmt.Println("链接至青云client error", err.Error())
		return nil, err
	}

	// please modify following params to send a test request
	ZONE := zone
	STATUS := []string{"active", "stopped"}
	LIMIT := 299
	OFFSET := 200

	resp, err := client.Boss2DescribeLoadBalancers(ZONE, STATUS, LIMIT, OFFSET)
	if err != nil {
		if err != nil {
			logger.Error("青云 LoadBalancers 获取 error", err.Error())
			fmt.Println("青云 LoadBalancers 获取 error", err.Error())
			return nil, err
		}
	}
	return models.ParseLoadBalancers([]byte(resp)), nil
}
