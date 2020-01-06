package pv

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	PVScrapperTypes = "PVScrapper"
	PVResourceTypes = "PV"
)

type PVScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewPVScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *PVScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
		kubernetesLister: lister,
	}

	return &PVScrapper{
		kubeAccessor:   ka,
		CommonScrapper: common.NewCommonScrapper(dc),
	}
}

func (w *PVScrapper) GetScrapperTypes() string {
	return PVScrapperTypes
}

func (w *PVScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, namespace)

}

func (w *PVScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	PVInfos, err := w.kubeAccessor.getPVs(namespace)
	if err != nil {
		return nil, err
	}

	return PVInfos, nil
}
