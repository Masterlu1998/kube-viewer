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
	ReplicaSetDetailScrapperTypes = "ReplicaSetDetailScrapper"
)

// Node detail scrapper
type ReplicaSetDetailScrapper struct {
	*common.CommonScrapper
}

func NewReplicaSetDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *ReplicaSetDetailScrapper {
	return &ReplicaSetDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *ReplicaSetDetailScrapper) GetScrapperTypes() string {
	return ReplicaSetDetailScrapperTypes
}

func (w *ReplicaSetDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *ReplicaSetDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	repicaSet, err := w.CommonScrapper.KubernetesLister.ReplicaSetLister.ReplicaSets(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(repicaSet)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
