package namespace

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
}

func (ka *kubeAccessor) getNamespaces() ([]string, error) {
	namespaces, err := ka.kubernetesClient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespacesInfos []string
	for _, ns := range namespaces.Items {
		if ns.Status.Phase == corev1.NamespaceActive {
			namespacesInfos = append(namespacesInfos, ns.Name)
		}
	}

	return namespacesInfos, nil
}
