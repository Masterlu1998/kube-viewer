package namespace

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	NamespaceScrapperTypes = "NamespaceScrapper"
	NamespaceResourceTypes = "Namespace"
)

type NamespaceScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewNamespaceScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *NamespaceScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
		kubernetesLister: lister,
	}

	return &NamespaceScrapper{
		kubeAccessor:   ka,
		CommonScrapper: common.NewCommonScrapper(dc),
	}
}

func (w *NamespaceScrapper) GetScrapperTypes() string {
	return NamespaceScrapperTypes
}

func (w *NamespaceScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, "")

}

func (w *NamespaceScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	namespaceInfos, err := w.kubeAccessor.getNamespaces()
	if err != nil {
		return nil, err
	}

	return namespaceInfos, nil
}
