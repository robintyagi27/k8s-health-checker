package health

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	NodeReadyGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "k8s_nodes_ready_total",
		Help: "Number of ready Kubernetes nodes",
	})
	PodRunningGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "k8s_pods_running_total",
		Help: "Number of running Kubernetes pods",
	})
	PodFailedGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "k8s_pods_failed_total",
		Help: "Number of failed Kubernetes pods",
	})
)

func init() {
	prometheus.MustRegister(NodeReadyGauge, PodRunningGauge, PodFailedGauge)
}

func StartMetricsServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":"+port, nil)
}
