package secret

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

func (ka *kubeAccessor) getSecrets(namespace string) ([]Info, error) {
	secrets, err := ka.kubernetesLister.SecretLister.Secrets(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var secretInfos []Info
	for _, s := range secrets {
		secretInfos = append(secretInfos, Info{
			Name:       s.Name,
			Namespace:  s.Namespace,
			CreateTime: s.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
		})
	}

	sort.Slice(secretInfos, func(left, right int) bool {
		return secretInfos[left].Name > secretInfos[right].Name
	})

	return secretInfos, nil
}
