package namespace

import (
	"sort"

	"github.com/Masterlu1998/kube-viewer/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
	kubernetesLister *kube.KubeLister
}

func (ka *kubeAccessor) getNamespaces() ([]string, error) {
	namespaces, err := ka.kubernetesLister.NamespaceLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var namespacesInfos []string
	for _, ns := range namespaces {
		if ns.Status.Phase == corev1.NamespaceActive {
			namespacesInfos = append(namespacesInfos, ns.Name)
		}
	}

	sort.Strings(namespacesInfos)

	return namespacesInfos, nil
}
