package servser

import (
	bossclient "exporter-demo/client"
	"exporter-demo/exporter_config"
	"exporter-demo/logger"
	"exporter-demo/models"
	"fmt"
)

func Wjdev_Listeners() *models.ListenerSetResponse {

	host, protocol, access_key, secret_key, zone, port, ignore_verify := exporter_config.ConfigBoss()

	// init client
	client, err := bossclient.NewClient(host, port, protocol, access_key, secret_key, ignore_verify)
	if err != nil {
		logger.Error("链接至青云client error", err.Error())
		fmt.Println("链接至青云client error", err.Error())
		return nil
	}

	// please modify following params to send a test request
	ZONE := zone
	LIMIT := 99999
	OFFSET := 0

	resp, err := client.DescribeLoadBalancerListeners(ZONE, LIMIT, OFFSET)
	//fmt.Println(resp)
	if err != nil {
		fmt.Println(err.Error())
	}
	return models.ParseListeners([]byte(resp))

	//return nil
}
