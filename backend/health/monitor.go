package health

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

func StartMetricsServer(port string) {
	prometheus.MustRegister(NodeReadyGauge, PodRunningGauge, PodFailedGauge)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Printf("Starting metrics server on :%s/metrics ...", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()
}

type ClusterStats struct {
	NodeCount    int
	ReadyNodes   int
	TotalPods    int
	RunningPods  int
	FailedPods   int
}

func CheckClusterHealth(clientset *kubernetes.Clientset) (*ClusterStats, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	stats := &ClusterStats{
		NodeCount: len(nodes.Items),
		TotalPods: len(pods.Items),
	}

	for _, node := range nodes.Items {
		for _, cond := range node.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				stats.ReadyNodes++
			}
		}
	}

	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case "Running":
			stats.RunningPods++
		case "Failed":
			stats.FailedPods++
		}
	}

	return stats, nil
}