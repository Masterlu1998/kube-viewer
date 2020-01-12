package secret

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
	SecretScrapperTypes = "SecretScrapper"
	SecretResourceTypes = "Secret"
)

type SecretScrapper struct {
	*common.CommonScrapper
}

func NewSecretScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *SecretScrapper {

	return &SecretScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *SecretScrapper) GetScrapperTypes() string {
	return SecretScrapperTypes
}

func (w *SecretScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *SecretScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	secretInfos, err := w.getSecrets(listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return secretInfos, nil
}

type Info struct {
	Name       string
	Namespace  string
	CreateTime string
}

func (w *SecretScrapper) getSecrets(namespace string) ([]Info, error) {
	secrets, err := w.KubernetesLister.SecretLister.Secrets(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var secretInfos []Info
	for _, s := range secrets {
		secretInfos = append(secretInfos, Info{
			Name:       s.Name,
			Namespace:  s.Namespace,
			CreateTime: s.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
		})
	}

	sort.Slice(secretInfos, func(left, right int) bool {
		return secretInfos[left].Name > secretInfos[right].Name
	})

	return secretInfos, nil
}
