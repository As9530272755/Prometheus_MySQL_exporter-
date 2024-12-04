package collectors

import (
	"exporter-demo/models"
	"exporter-demo/servser"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
	"time"
)

var (
	LBStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "qingcloud_loadbalancer_status",
			Help: "Status of the QingCloud LoadBalancer",
		},
		[]string{"zone", "lb_id", "lb_name", "lb_status", "lb_ip"},
	)

	LBInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "qingcloud_loadbalancer_info",
			Help: "Info of the QingCloud LoadBalancer",
		},
		[]string{"zone", "lb_id", "lb_name", "lb_status", "lb_ip", "rsyslog", "securityGroup", "privateIP"})

	LBMaxConnection = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "qingcloud_loadbalancer_max_conn",
			Help: "Max Connection of the QingCloud LoadBalancer",
		},
		[]string{"zone", "lb_id", "lb_name", "lb_status", "lb_ip"})

	//ListenerTotalConnections = prometheus.NewGaugeVec(
	//	prometheus.GaugeOpts{
	//		Name: "qingcloud_listener_total_connections",
	//		Help: "Max Connection of the QingCloud LoadBalancer",
	//	},
	//	[]string{"zone", "lb_id", "listener_name", "istener_id"})
	//[]string{"istener_id"})
)

func init() {
	prometheus.MustRegister(LBStatus)
	prometheus.MustRegister(LBInfo)
	prometheus.MustRegister(LBMaxConnection)
	//prometheus.MustRegister(ListenerTotalConnections)
}

// 由于青云 LB 每次最多只能获取 100 个，所以这里将不同分页进行数据整合
func FetchLoadBalancers() []models.LoadbalancerSet {

	lbs1, err := servser.Wjdev_lb_set1()
	if err != nil {
		// handle error appropriately
		return nil
	}
	lbs2, err := servser.Wjdev_lb_set2()
	if err != nil {
		// handle error appropriately
		return nil
	}
	lbs3, err := servser.Wjdev_lb_set3()
	if err != nil {
		// handle error appropriately
		return nil
	}

	// 合并获取所有 lb
	allLoadBalancers := lbs1.LoadbalancerSets
	allLoadBalancers = append(allLoadBalancers, lbs2.LoadbalancerSets...)
	allLoadBalancers = append(allLoadBalancers, lbs3.LoadbalancerSets...)

	return allLoadBalancers
}

//func FetchListeners() []models.Listener {
//	Listener := servser.Wjdev_Listeners()
//	return Listener.LoadbalancerListenerSet
//}

func LBMetrics(wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()

	// 10s 更新一次 API,创建了一个定时器 ticker，每隔 10 秒触发一次
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var lastLoadBalancers []models.LoadbalancerSet

	for {
		select {
		case <-ticker.C:
			loadBalancers := FetchLoadBalancers()

			//listeners := FetchListeners()

			//数据无变化，跳过更新指标，用久数据和新数据做对比，并且 equal 返回得数据为真才跳过更新
			if len(lastLoadBalancers) == len(loadBalancers) && equal(lastLoadBalancers, loadBalancers) {
				continue
			}

			// 重置 exporter 获取得 metrics
			LBStatus.Reset()
			LBInfo.Reset()
			LBMaxConnection.Reset()
			//ListenerTotalConnections.Reset()
			// 对状态做判断得 Metrics
			for _, lb := range loadBalancers {

				status := 0.0
				if lb.Status == "active" {
					status = 1.0
				} else {
					status = 0.0
				}

				// 获取 LBIP
				lbip := ""

				for _, cluster := range lb.Cluster {
					if cluster.EipName != "" && cluster.EipAddr != "" {
						lbip = cluster.EipAddr
					}
				}

				// 获取安全组
				securityGroup := ""
				for _, SecurityGroup := range lb.SecurityGroups {
					if SecurityGroup.GroupName != "" {
						securityGroup = SecurityGroup.GroupName
					}
				}

				// 获取私有IP集
				privateIP := strings.Join(lb.PrivateIPs, ",")
				// LB状态做监控
				LBStatus.WithLabelValues(lb.ZoneID, lb.LoadbalancerID, lb.LoadbalancerName, lb.Status, lbip).Set(status)
				// LB信息监控
				LBInfo.WithLabelValues(lb.ZoneID, lb.LoadbalancerID, lb.LoadbalancerName, lb.Status, lbip, lb.Rsyslog, securityGroup, privateIP).Set(1)

				maxConn := lb.LoadbalancerType
				switch maxConn {
				case 0:
					LBMaxConnection.WithLabelValues(lb.ZoneID, lb.LoadbalancerID, lb.LoadbalancerName, lb.Status, lbip).Set(5000)
				case 1:
					LBMaxConnection.WithLabelValues(lb.ZoneID, lb.LoadbalancerID, lb.LoadbalancerName, lb.Status, lbip).Set(20000)
				case 2:
					LBMaxConnection.WithLabelValues(lb.ZoneID, lb.LoadbalancerID, lb.LoadbalancerName, lb.Status, lbip).Set(40000)
				case 3:
					LBMaxConnection.WithLabelValues(lb.ZoneID, lb.LoadbalancerID, lb.LoadbalancerName, lb.Status, lbip).Set(100000)
				}

				lastLoadBalancers = loadBalancers
			}

		case <-stopChan:
			return
		}
	}
}

// 用于判断 LB 数据是否需要更新
func equal(a, b []models.LoadbalancerSet) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		// 判断 LB 状态是否变化，如果有变化也需要更新数据
		if a[i].Status != b[i].Status {
			return false
		} else if a[i].LoadbalancerType != b[i].LoadbalancerType { // 判断 LB 连接数是否变化，如果有变化也需要更新数据
			return false
		}
	}
	return true
}
