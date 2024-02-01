package cmd

import (
	"context"
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func main() {
	// Konfiguracja logowania i flag
	klog.InitFlags(nil)
	flag.Parse()

	// Ładowanie konfiguracji klienta Kubernetes
	config, err := getClientConfig()
	if err != nil {
		klog.Fatalf("Błąd podczas ładowania konfiguracji klienta: %s", err.Error())
	}

	// Tworzenie klienta Kubernetes
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Błąd podczas tworzenia klienta Kubernetes: %s", err.Error())
	}

	// Inicjalizacja niestandardowego planisty
	customScheduler := scheduler.NewCustomScheduler(clientset)
	klog.Info("Niestandardowy planista został zainicjalizowany")

	// Uruchomienie pętli obserwującej
	customScheduler.Run(context.Background())
}

func getClientConfig() (*rest.Config, error) {
	kubeconfig := flag.String("kubeconfig", "", "Ścieżka do pliku kubeconfig (opcjonalnie)")
	if *kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", *kubeconfig)
	}
	return rest.InClusterConfig()
}
