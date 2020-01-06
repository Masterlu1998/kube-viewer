package node

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	NodeScrapperTypes = "NodeScrapper"
	NodeResourceTypes = "Node"
)

type NodeScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewNodeScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *NodeScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
		kubernetesLister: lister,
	}

	return &NodeScrapper{
		kubeAccessor:   ka,
		CommonScrapper: common.NewCommonScrapper(dc),
	}
}

func (w *NodeScrapper) GetScrapperTypes() string {
	return NodeScrapperTypes
}

func (w *NodeScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, namespace)

}

func (w *NodeScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	nodeInfos, err := w.kubeAccessor.getNodes(namespace)
	if err != nil {
		return nil, err
	}

	return nodeInfos, nil
}
