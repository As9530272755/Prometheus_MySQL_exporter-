package main

import (
	"exporter-demo/collectors"
	"exporter-demo/exporter_config"
	"exporter-demo/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	//collectors.FetchListeners()
	logger.InitLog("error.log", true)
	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	wg.Add(1)
	go collectors.LBMetrics(&wg, stopChan)

	go collectors.ListenerMetrics(&wg, stopChan)

	port := exporter_config.LbConfig()
	addr := ":" + port

	//internal.WJLoadbalancerMetrics()

	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(addr, nil)

	//等待中断信号以正常关闭应用程序，接收信号(通常是 Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	//向goroutine发出停止信号
	//close(stopChan)
	<-stopChan

	//等待所有goroutines完成

	wg.Wait()
	//servser.Bjdev_lb_set2()
	//servser.Bjdev_lb_set3()

}
