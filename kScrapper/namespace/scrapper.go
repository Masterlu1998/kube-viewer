package namespace

import (
	"context"
	"sort"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

const (
	NamespaceScrapperTypes = "NamespaceScrapper"
	NamespaceResourceTypes = "Namespace"
)

type NamespaceScrapper struct {
	*common.CommonScrapper
}

func NewNamespaceScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *NamespaceScrapper {
	return &NamespaceScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *NamespaceScrapper) GetScrapperTypes() string {
	return NamespaceScrapperTypes
}

func (w *NamespaceScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, "")

}

func (w *NamespaceScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	namespaceInfos, err := w.getNamespaces()
	if err != nil {
		return nil, err
	}

	return namespaceInfos, nil
}

func (w *NamespaceScrapper) getNamespaces() ([]string, error) {
	namespaces, err := w.KubernetesLister.NamespaceLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var namespacesInfos []string
	for _, ns := range namespaces {
		if ns.Status.Phase == corev1.NamespaceActive {
			namespacesInfos = append(namespacesInfos, ns.Name)
		}
	}

	sort.Strings(namespacesInfos)

	return namespacesInfos, nil
}
