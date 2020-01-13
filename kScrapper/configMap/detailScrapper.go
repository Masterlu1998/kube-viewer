package configMap

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	ConfigMapDetailScrapperTypes = "ConfigMapDetailScrapper"
)

// Node detail scrapper
type ConfigMapDetailScrapper struct {
	*common.CommonScrapper
}

func NewConfigMapDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *ConfigMapDetailScrapper {
	return &ConfigMapDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *ConfigMapDetailScrapper) GetScrapperTypes() string {
	return ConfigMapDetailScrapperTypes
}

func (w *ConfigMapDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *ConfigMapDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	configMap, err := w.CommonScrapper.KubernetesLister.ConfigMapLister.ConfigMaps(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(configMap)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
