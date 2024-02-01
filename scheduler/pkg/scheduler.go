package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // metav1 alias defined here
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	influxToken  = "your-influxdb-token"
	influxOrg    = "your-influxdb-org"
	influxBucket = "your-influxdb-bucket"
	influxURL    = "http://localhost:8086"
)

type CustomScheduler struct {
	clientset *kubernetes.Clientset
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
}

func NewCustomScheduler(clientset *kubernetes.Clientset) *CustomScheduler {
	scheduler := &CustomScheduler{
		clientset: clientset,
		queue:     workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	podInformer := cache.NewSharedIndexInformer(
		cache.NewListWatchFromClient(scheduler.clientset.CoreV1().RESTClient(), "pods", v1.NamespaceAll, fields.Everything()),
		&v1.Pod{},
		0,
		cache.Indexers{},
	)

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: scheduler.enqueuePod,
	})

	scheduler.informer = podInformer

	return scheduler
}

func (s *CustomScheduler) Run(ctx context.Context) {
	defer s.queue.ShutDown()

	klog.Info("Starting custom scheduler")
	go s.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), s.informer.HasSynced) {
		klog.Error("Error syncing cache")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.scheduleNextPod(ctx)
		}
	}
}

func (s *CustomScheduler) enqueuePod(obj interface{}) {
	klog.Info("Enqueuing pod")
	s.queue.Add(obj)
}

func (s *CustomScheduler) scheduleNextPod(ctx context.Context) {
	klog.Info("Attempting to schedule next pod")
	obj, shutdown := s.queue.Get()

	if shutdown {
		return
	}

	defer s.queue.Done(obj)

	pod, ok := obj.(*v1.Pod)
	if !ok {
		klog.Error("Expected Pod type")
		return
	}

	selectedNode := "node-name" // Replace with scheduling logic

	err := s.schedulePod(pod, selectedNode)
	if err != nil {
		klog.Errorf("Error scheduling pod: %v", err)
		return
	}

	klog.Infof("Pod %s/%s scheduled to node %s", pod.Namespace, pod.Name, selectedNode)
}

func (s *CustomScheduler) schedulePod(pod *v1.Pod, nodeName string) error {
	patchBytes := []byte(fmt.Sprintf(`{"spec":{"nodeName":"%s"}}`, nodeName))
	_, err := s.clientset.CoreV1().Pods(pod.Namespace).Patch(context.TODO(), pod.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("error assigning pod %s to node %s: %v", pod.Name, nodeName, err)
	}

	return nil
}

func (s *CustomScheduler) getAvailableNodes(ctx context.Context) ([]v1.Node, error) {
	nodeList, err := s.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching node list: %v", err)
	}
	return nodeList.Items, nil // Make sure it returns []v1.Node, not []*v1.Node
}

func (s *CustomScheduler) selectBestNode(pod *v1.Pod, nodes []*v1.Node) string {
	var bestNode *v1.Node
	bestScore := -1.0

	for _, node := range nodes {
		score := s.evaluateNode(node, pod)

		if score > bestScore {
			bestNode = node
			bestScore = score
		}
	}

	if bestNode != nil {
		return bestNode.Name
	}

	return ""
}

func (s *CustomScheduler) evaluateNode(node *v1.Node, pod *v1.Pod) float64 {
	totalMemory := node.Status.Allocatable.Memory().Value()
	usedMemory := int64(0) // Implement memory usage retrieval

	predictedLoad, err := s.getPredictedLoad(node.Name)
	if err != nil {
		klog.Errorf("Failed to get predicted load for node %s: %v", node.Name, err)
		predictedLoad = 0.5 // Default load value for demonstration
	}

	score := float64(totalMemory-usedMemory) * (1.0 - predictedLoad)
	return score
}

func (s *CustomScheduler) getPredictedLoad(nodeName string) (float64, error) {
	// Utwórz nowego klienta InfluxDB
	client := influxdb2.NewClient(influxURL, influxToken)
	defer client.Close()

	// Pobierz interfejs API zapytań
	queryAPI := client.QueryAPI(influxOrg)

	// Przygotuj zapytanie do InfluxDB
	query := fmt.Sprintf(`from(bucket:"%s")|> range(start: -1h)|> filter(fn:(r) => r._measurement == "node_load" and r.node == "%s")`, influxBucket, nodeName)

	// Wykonaj zapytanie
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return 0, fmt.Errorf("error querying influxdb for node %s: %v", nodeName, err)
	}

	// Zmienna na przechowanie prognozowanego obciążenia
	var predictedLoad float64

	// Iteruj przez wyniki zapytania
	for result.Next() {
		// Pobierz wartość obciążenia z bieżącego rekordu
		predictedLoad = result.Record().ValueByKey("_value").(float64)
	}

	if result.Err() {
		return 0, fmt.Errorf("error during query iteration for node %s: %v", nodeName, result.Err())
	}

	return predictedLoad, nil
}
