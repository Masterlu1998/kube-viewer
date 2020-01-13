package pvc

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	PVCDetailScrapperTypes = "PVCDetailScrapper"
)

// Node detail scrapper
type PVCDetailScrapper struct {
	*common.CommonScrapper
}

func NewPVCDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *PVCDetailScrapper {
	return &PVCDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *PVCDetailScrapper) GetScrapperTypes() string {
	return PVCDetailScrapperTypes
}

func (w *PVCDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *PVCDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	pvc, err := w.CommonScrapper.KubernetesLister.PVCLister.PersistentVolumeClaims(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(pvc)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
