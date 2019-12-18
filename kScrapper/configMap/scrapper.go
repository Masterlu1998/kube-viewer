package configMap

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	ConfigMapScrapperTypes = "ConfigMapScrapper"
	ConfigMapResourceTypes = "ConfigMap"
)

type ConfigMapScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewConfigMapScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *ConfigMapScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
		kubernetesLister: lister,
	}

	return &ConfigMapScrapper{
		kubeAccessor:   ka,
		CommonScrapper: common.NewCommonScrapper(dc),
	}
}

func (w *ConfigMapScrapper) GetScrapperTypes() string {
	return ConfigMapScrapperTypes
}

func (w *ConfigMapScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, namespace)

}

func (w *ConfigMapScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	configMapInfos, err := w.kubeAccessor.getConfigMaps(namespace)
	if err != nil {
		return nil, err
	}

	return configMapInfos, nil
}
