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
	DaemonSetDetailScrapperTypes = "DaemonSetDetailScrapper"
)

// Node detail scrapper
type DaemonSetDetailScrapper struct {
	*common.CommonScrapper
}

func NewDaemonSetDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *DaemonSetDetailScrapper {
	return &DaemonSetDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *DaemonSetDetailScrapper) GetScrapperTypes() string {
	return DaemonSetDetailScrapperTypes
}

func (w *DaemonSetDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *DaemonSetDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	daemonSet, err := w.CommonScrapper.KubernetesLister.DaemonSetLister.DaemonSets(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(daemonSet)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
