package workload

import (
	"context"
	"errors"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	ReplicaSetScrapperTypes = "ReplicaSetScrapper"
	ReplicaSetResourceTypes = "ReplicaSet"
)

type ReplicaSetScrapper struct {
	*common.CommonScrapper
}

func NewReplicaSetListScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *ReplicaSetScrapper {
	return &ReplicaSetScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (r *ReplicaSetScrapper) GetScrapperTypes() string {
	return ReplicaSetScrapperTypes
}

func (r *ReplicaSetScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	r.CommonScrapper.ScrapeDataIntoChWithSource(ctx, r.scrapeDataIntoCh, args)
}

func (r *ReplicaSetScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	replicaSets, err := getWorkloads(r.KubernetesClient, r.KubernetesLister, ReplicaSetResourceTypes, listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return replicaSets, nil
}
