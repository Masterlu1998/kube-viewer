package node

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	NodeDetailScrapperTypes = "NodeDetailScrapper"
)

// Node detail scrapper
type NodeDetailScrapper struct {
	*common.CommonScrapper
}

func NewNodeDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *NodeDetailScrapper {
	return &NodeDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *NodeDetailScrapper) GetScrapperTypes() string {
	return NodeDetailScrapperTypes
}

func (w *NodeDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *NodeDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	node, err := w.CommonScrapper.KubernetesLister.NodeLister.Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(node)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
