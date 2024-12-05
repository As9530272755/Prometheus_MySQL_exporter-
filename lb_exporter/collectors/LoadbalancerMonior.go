package collectors

import (
	"exporter-demo/logger"
	"exporter-demo/models"
	"exporter-demo/servser"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ListenerAccumulatedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "qingcloud_listener_accumulated_connections",
			Help: "Accumulated Connections of the QingCloud LoadBalancer",
		},
		[]string{"zone", "lb_id", "listener_name", "istener_id"})
	ListenerMaxConcurrencys = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "qingcloud_listener_max_concurrency",
			Help: "Max Concurrency of the QingCloud LoadBalancer",
		},
		[]string{"zone", "lb_id", "listener_name", "istener_id"})
	ListenerAverageConcurrency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "qingcloud_listener_average_concurrency",
			Help: "Average Concurrency of the QingCloud LoadBalancer",
		},
		[]string{"zone", "lb_id", "listener_name", "istener_id"})
	ListenerConcurrentLimit = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "qingcloud_listener_concurrent_limit",
			Help: "Concurrent Limit of the QingCloud LoadBalancer",
		},
		[]string{"zone", "lb_id", "listener_name", "istener_id"})
)

func init() {
	prometheus.MustRegister(ListenerAccumulatedConnections)
	prometheus.MustRegister(ListenerMaxConcurrencys)
	prometheus.MustRegister(ListenerAverageConcurrency)
	prometheus.MustRegister(ListenerConcurrentLimit)
}

// 获取所有的 Listener 信息
func FetchListeners() []models.Listener {
	Listener := servser.Wjdev_Listeners()
	return Listener.LoadbalancerListenerSet
}

func ListenerMetrics(wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()

	// 10M 更新一次 API,创建了一个定时器 ticker，每隔 10 M触发一次
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			defer func() {
				if r := recover(); r != nil {
					logger.Error("Recovered from panic:", r)
					fmt.Println("Recovered from panic:", r)
				}
			}()
			listeners := FetchListeners()
			ListenerAccumulatedConnections.Reset()
			ListenerConcurrentLimit.Reset()
			ListenerMaxConcurrencys.Reset()
			ListenerAverageConcurrency.Reset()

			// 采集所有监听器，并且获取到对应的数据
			for _, listener := range listeners {
				// 使用 map 去重并处理 NA
				uniqueDataMap := make(map[string]bool)
				var uniqueData []interface{}
				// 总连接数
				AccumulatedConnections := []string{}
				// 最大并发连接数
				MaxConcurrency := []string{}
				// 并发连接数上限
				ConcurrentLimit := []string{}
				// 平均值
				AverageConcurrency := []string{}

				loadbalancerMonior, err := servser.LoadBalancerMonitor()
				//fmt.Println(loadbalancerMonior)
				if err != nil {
					logger.Error("Error loadbalancerMonior: %v", err)
					fmt.Println("Error loadbalancerMonior: %v", err)
				}

				if len(loadbalancerMonior.MeterSet) > 0 && len(loadbalancerMonior.MeterSet[0].DataSet) > 0 && len(loadbalancerMonior.MeterSet[0].DataSet[0].Data) > 0 {
					dataSet := loadbalancerMonior.MeterSet[0].DataSet[0].Data
					// 提取 DataSet
					for _, item := range dataSet {
						strItem, ok := item.(string)
						if ok && strItem != "NA" && !uniqueDataMap[strItem] {
							uniqueData = append(uniqueData, strItem)
							uniqueDataMap[strItem] = true
						} else if numItem, ok := item.(float64); ok { // 处理数字类型
							numStrItem := fmt.Sprintf("%v", numItem)
							if !uniqueDataMap[numStrItem] {
								uniqueData = append(uniqueData, numItem)
								fmt.Println(uniqueData)
								uniqueDataMap[numStrItem] = true
							}
						}
					}

					AccumulatedConnections = append(AccumulatedConnections, strings.Split(uniqueData[0].(string), `|`)[0])

					MaxConcurrency = append(MaxConcurrency, strings.Split(uniqueData[0].(string), `|`)[1])

					AverageConcurrency = append(AverageConcurrency, strings.Split(uniqueData[0].(string), `|`)[2])

					ConcurrentLimit = append(ConcurrentLimit, strings.Split(uniqueData[0].(string), `|`)[len(strings.Split(uniqueData[0].(string), `|`))-1])

				} else {
					defer func() {
						if r := recover(); r != nil {
							logger.Error("Recovered from panic: %v", r)
							fmt.Println("Recovered from panic:", r)
							// 这里可以选择继续执行程序的其他部分，而不是退出
						}
					}()
				}

				accumulatedConnectionsConnections, _ := strconv.ParseFloat(AccumulatedConnections[0], 64)
				maxConcurrency, _ := strconv.ParseFloat(MaxConcurrency[0], 64)
				averageConcurrency, _ := strconv.ParseFloat(AverageConcurrency[0], 64)
				concurrentLimit, _ := strconv.ParseFloat(ConcurrentLimit[0], 64)

				ListenerAccumulatedConnections.WithLabelValues(listener.ZoneID, listener.LoadbalancerID, listener.ListenerName, loadbalancerMonior.ResourceID).Set(accumulatedConnectionsConnections)
				ListenerMaxConcurrencys.WithLabelValues(listener.ZoneID, listener.LoadbalancerID, listener.ListenerName, loadbalancerMonior.ResourceID).Set(maxConcurrency)
				ListenerAverageConcurrency.WithLabelValues(listener.ZoneID, listener.LoadbalancerID, listener.ListenerName, loadbalancerMonior.ResourceID).Set(averageConcurrency)
				ListenerConcurrentLimit.WithLabelValues(listener.ZoneID, listener.LoadbalancerID, listener.ListenerName, loadbalancerMonior.ResourceID).Set(concurrentLimit)

				//ListenerTotalConnections.WithLabelValues(loadbalancerMonior.ResourceID).Set(totalConnections)
			}

			//}
		case <-stopChan:
			return
		}
	}
}
