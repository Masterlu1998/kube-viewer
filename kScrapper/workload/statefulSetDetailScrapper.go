package workload

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	StatefulSetDetailScrapperTypes = "StatefulSetDetailScrapper"
)

// Node detail scrapper
type StatefulSetDetailScrapper struct {
	*common.CommonScrapper
}

func NewStatefulSetDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *StatefulSetDetailScrapper {
	return &StatefulSetDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *StatefulSetDetailScrapper) GetScrapperTypes() string {
	return StatefulSetDetailScrapperTypes
}

func (w *StatefulSetDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *StatefulSetDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	statefulSet, err := w.CommonScrapper.KubernetesLister.StatefulSetLister.StatefulSets(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(statefulSet)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
