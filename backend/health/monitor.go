package health

import (
	"context"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ClusterHealth struct {
	NodeCount     int
	ReadyNodes    int
	NotReadyNodes int
	TotalPods     int
	RunningPods   int
	FailedPods    int
	PendingPods   int
}

func CheckClusterHealth(clientset *kubernetes.Clientset) (*ClusterHealth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	health := &ClusterHealth{}

	// --- Check nodes ---
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %v", err)
	}
	health.NodeCount = len(nodes.Items)
	for _, node := range nodes.Items {
		for _, cond := range node.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				health.ReadyNodes++
			}
		}
	}
	health.NotReadyNodes = health.NodeCount - health.ReadyNodes

	// --- Check pods ---
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}
	health.TotalPods = len(pods.Items)
	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case "Running":
			health.RunningPods++
		case "Failed":
			health.FailedPods++
		case "Pending":
			health.PendingPods++
		}
	}

	log.Printf("[HealthCheck] Nodes: %d ready/%d total | Pods: %d running/%d failed\n",
		health.ReadyNodes, health.NodeCount, health.RunningPods, health.TotalPods)

	return health, nil
}
