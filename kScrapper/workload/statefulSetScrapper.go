package workload

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	StatefulSetScrapperTypes = "StatefulSetScrapper"
	StatefulSetResourceTypes = "StatefulSet"
)

type StatefulSetScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewStatefulSetScrapper(lister *kube.KubeLister, client *kubernetes.Clientset) *StatefulSetScrapper {
	return &StatefulSetScrapper{
		kubeAccessor:   generateKubeAccessor(lister, client),
		CommonScrapper: common.NewCommonScrapper(),
	}
}

func (w *StatefulSetScrapper) GetScrapperTypes() string {
	return StatefulSetScrapperTypes
}

func (w *StatefulSetScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, namespace)
}

func (w *StatefulSetScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	statefulSets, err := w.kubeAccessor.getWorkloads(StatefulSetResourceTypes, namespace)
	if err != nil {
		return nil, err
	}

	return statefulSets, nil
}
