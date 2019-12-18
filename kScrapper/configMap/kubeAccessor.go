package configMap

import (
	"sort"
	"time"

	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	Name       string
	Namespace  string
	CreateTime string
}

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
	kubernetesLister *kube.KubeLister
}

func (ka *kubeAccessor) getConfigMaps(namespace string) ([]Info, error) {
	configMaps, err := ka.kubernetesLister.ConfigMapLister.ConfigMaps(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var configMapInfos []Info
	for _, s := range configMaps {
		configMapInfos = append(configMapInfos, Info{
			Name:       s.Name,
			Namespace:  s.Namespace,
			CreateTime: s.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
		})
	}

	sort.Slice(configMapInfos, func(left, right int) bool {
		return configMapInfos[left].Name > configMapInfos[right].Name
	})

	return configMapInfos, nil
}
