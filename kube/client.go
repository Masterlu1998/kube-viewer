package kube

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const kubeConfigFilePath = "./kube/admin.conf"

func GetKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func getConfig() (*rest.Config, error) {
	if _, err := os.Stat(kubeConfigFilePath); err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigFilePath)
	if err != nil {
		return nil, err
	}

	return config, nil
}
