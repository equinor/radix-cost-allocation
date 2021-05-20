package utils

import (
	"errors"
	"os"

	radixclientset "github.com/equinor/radix-operator/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubernetesClients returns clients for accessing k8s and Radix resources in a cluster
func GetKubernetesClients() (kubernetes.Interface, radixclientset.Interface, error) {
	config, err := getClientConfig()
	if err != nil {
		return nil, nil, err
	}

	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	radixclient, err := radixclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return kubeclient, radixclient, nil
}

func getClientConfig() (*rest.Config, error) {
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"
	if config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath); err == nil {
		return config, nil
	}

	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}

	return nil, errors.New("unable to get kubernetes config")
}
