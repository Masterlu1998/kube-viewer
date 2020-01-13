package workload

import (
	"context"
	"errors"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	StatefulSetListScrapperTypes = "StatefulSetScrapper"
	StatefulSetResourceTypes     = "StatefulSet"
)

type StatefulSetScrapper struct {
	*common.CommonScrapper
}

func NewStatefulSetListScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *StatefulSetScrapper {
	return &StatefulSetScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *StatefulSetScrapper) GetScrapperTypes() string {
	return StatefulSetListScrapperTypes
}

func (w *StatefulSetScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)
}

func (w *StatefulSetScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	statefulSets, err := getWorkloads(w.KubernetesClient, w.KubernetesLister, StatefulSetResourceTypes, listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return statefulSets, nil
}
