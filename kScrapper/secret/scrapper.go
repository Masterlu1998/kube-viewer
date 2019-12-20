package secret

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	SecretScrapperTypes = "SecretScrapper"
	SecretResourceTypes = "Secret"
)

type SecretScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewSecretScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *SecretScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
		kubernetesLister: lister,
	}

	return &SecretScrapper{
		kubeAccessor:   ka,
		CommonScrapper: common.NewCommonScrapper(dc),
	}
}

func (w *SecretScrapper) GetScrapperTypes() string {
	return SecretScrapperTypes
}

func (w *SecretScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, namespace)

}

func (w *SecretScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	secretInfos, err := w.kubeAccessor.getSecrets(namespace)
	if err != nil {
		return nil, err
	}

	return secretInfos, nil
}
