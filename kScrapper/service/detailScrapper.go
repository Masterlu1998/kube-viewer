package service

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	ServiceDetailScrapperTypes = "ServiceDetailScrapper"
)

// Node detail scrapper
type ServiceDetailScrapper struct {
	*common.CommonScrapper
}

func NewServiceDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *ServiceDetailScrapper {
	return &ServiceDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *ServiceDetailScrapper) GetScrapperTypes() string {
	return ServiceDetailScrapperTypes
}

func (w *ServiceDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *ServiceDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	service, err := w.CommonScrapper.KubernetesLister.ServiceLister.Services(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(service)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
