package collectors

import (
	"exporter-demo/logger"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"strconv"
	"time"
)

var (
	PortStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "LoadBalancersAndListener_Port_Status",
			Help: "LoadBalancersAndListener status (1 for reachale, 0 for net reachable)",
		},
		[]string{"lb_ip", "lb_id", "lb_name", "zone", "listener_id", "listener_name", "listener_port"},
	)
)

func init() {
	prometheus.MustRegister(PortStatusGauge)
}

func probeIPPort(ip, port string) bool {
	if ip == "" {
		return false // Skip if IP is empty
	}

	address := net.JoinHostPort(ip, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		logger.Error("Failed to connect to %s: %v\n", address, err)
		fmt.Println("Failed to connect to %s: %v\n", address, err)
		return false
	}

	fmt.Println("port ip", port, ip)
	defer conn.Close()
	return true
}

func ListenersPortMetrics(lb_ip, lb_id, lb_name, lb_zone, listener_id, listener_name string, listener_port int) {

	port := strconv.Itoa(listener_port)

	//for {
	//	conn, err := net.DialTimeout("tcp", net.JoinHostPort(lb_ip, port))
	//	if err != nil {
	//		PortStatusGauge.WithLabelValues(lb_ip, lb_id, lb_name, lb_zone, listener_id, listener_name, port).Set(0)
	//	} else {
	//		PortStatusGauge.WithLabelValues(lb_ip, lb_id, lb_name, lb_zone, listener_id, listener_name, port).Set(1)
	//		conn.Close()
	//	}
	//}

	success := probeIPPort(lb_ip, port)
	if success {
		PortStatusGauge.WithLabelValues(lb_ip, lb_id, lb_name, lb_zone, listener_id, listener_name, port).Set(1)
	} else {
		PortStatusGauge.WithLabelValues(lb_ip, lb_id, lb_name, lb_zone, listener_id, listener_name, port).Set(0)
	}

	// 每 10s 检查一次
	time.Sleep(10 * time.Second)

}
