package service

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	ServiceScrapperTypes = "ServiceScrapper"
	ServiceResourceTypes = "Namespace"
)

type ServiceScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewNamespaceScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *ServiceScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
		kubernetesLister: lister,
	}

	return &ServiceScrapper{
		kubeAccessor:   ka,
		CommonScrapper: common.NewCommonScrapper(dc),
	}
}

func (w *ServiceScrapper) GetScrapperTypes() string {
	return ServiceScrapperTypes
}

func (w *ServiceScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, "")

}

func (w *ServiceScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	namespaceInfos, err := w.kubeAccessor.getServices(namespace)
	if err != nil {
		return nil, err
	}

	return namespaceInfos, nil
}
