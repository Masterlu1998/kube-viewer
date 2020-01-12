package configMap

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

const (
	ConfigMapScrapperTypes = "ConfigMapScrapper"
	ConfigMapResourceTypes = "ConfigMap"
)

type ConfigMapScrapper struct {
	*common.CommonScrapper
}

func NewConfigMapScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *ConfigMapScrapper {
	return &ConfigMapScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *ConfigMapScrapper) GetScrapperTypes() string {
	return ConfigMapScrapperTypes
}

func (w *ConfigMapScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *ConfigMapScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	configMapInfos, err := w.getConfigMaps(listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return configMapInfos, nil
}

type Info struct {
	Name       string
	Namespace  string
	CreateTime string
}

func (w *ConfigMapScrapper) getConfigMaps(namespace string) ([]Info, error) {
	configMaps, err := w.KubernetesLister.ConfigMapLister.ConfigMaps(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var configMapInfos []Info
	for _, s := range configMaps {
		configMapInfos = append(configMapInfos, Info{
			Name:       s.Name,
			Namespace:  s.Namespace,
			CreateTime: s.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
		})
	}

	sort.Slice(configMapInfos, func(left, right int) bool {
		return configMapInfos[left].Name > configMapInfos[right].Name
	})

	return configMapInfos, nil
}
