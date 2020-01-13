package secret

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	SecretDetailScrapperTypes = "SecretDetailScrapper"
)

// Node detail scrapper
type SecretDetailScrapper struct {
	*common.CommonScrapper
}

func NewSecretDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *SecretDetailScrapper {
	return &SecretDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *SecretDetailScrapper) GetScrapperTypes() string {
	return SecretDetailScrapperTypes
}

func (w *SecretDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *SecretDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	secret, err := w.CommonScrapper.KubernetesLister.SecretLister.Secrets(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(secret)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
