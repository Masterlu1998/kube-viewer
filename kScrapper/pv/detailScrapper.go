package pv

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	PVDetailScrapperTypes = "PVDetailScrapper"
)

// Node detail scrapper
type PVDetailScrapper struct {
	*common.CommonScrapper
}

func NewPVDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *PVDetailScrapper {
	return &PVDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *PVDetailScrapper) GetScrapperTypes() string {
	return PVDetailScrapperTypes
}

func (w *PVDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *PVDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	pv, err := w.CommonScrapper.KubernetesLister.PVLister.Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(pv)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
