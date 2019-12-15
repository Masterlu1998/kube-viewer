package workload

import (
	"context"

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
	kubeAccessor *kubeAccessor
}

func NewReplicaSetScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *ReplicaSetScrapper {
	return &ReplicaSetScrapper{
		kubeAccessor:   generateKubeAccessor(lister, client),
		CommonScrapper: common.NewCommonScrapper(dc),
	}
}

func (r *ReplicaSetScrapper) GetScrapperTypes() string {
	return ReplicaSetScrapperTypes
}

func (r *ReplicaSetScrapper) StartScrapper(ctx context.Context, namespace string) {
	r.CommonScrapper.ScrapeDataIntoChWithSource(ctx, r.scrapeDataIntoCh, namespace)
}

func (r *ReplicaSetScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	replicaSets, err := r.kubeAccessor.getWorkloads(ReplicaSetResourceTypes, namespace)
	if err != nil {
		return nil, err
	}

	return replicaSets, nil
}
