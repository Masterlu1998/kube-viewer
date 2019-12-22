package pvc

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	PVCScrapperTypes = "PVCScrapper"
	PVCResourceTypes = "PVC"
)

type PVCScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewPVCScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *PVCScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
		kubernetesLister: lister,
	}

	return &PVCScrapper{
		kubeAccessor:   ka,
		CommonScrapper: common.NewCommonScrapper(dc),
	}
}

func (w *PVCScrapper) GetScrapperTypes() string {
	return PVCScrapperTypes
}

func (w *PVCScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, namespace)

}

func (w *PVCScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	configMapInfos, err := w.kubeAccessor.getPVCs(namespace)
	if err != nil {
		return nil, err
	}

	return configMapInfos, nil
}
