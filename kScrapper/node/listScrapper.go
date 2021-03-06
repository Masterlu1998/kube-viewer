package node

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

const (
	NodeListScrapperTypes  = "NodeListScrapper"
	NodeResourceTypes      = "Node"
	masterRoleSpecialLabel = "node-role.kubernetes.io/master"
)

// Node list scrapper
type NodeListScrapper struct {
	*common.CommonScrapper
}

func NewNodeListScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *NodeListScrapper {
	return &NodeListScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *NodeListScrapper) GetScrapperTypes() string {
	return NodeListScrapperTypes
}

func (w *NodeListScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *NodeListScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	nodeInfos, err := w.listNodes()
	if err != nil {
		return nil, err
	}

	return nodeInfos, nil
}

type Info struct {
	Name              string
	Status            string
	Roles             string
	Address           string
	OSImage           string
	CreationTimestamp string
}

func (w *NodeListScrapper) listNodes() ([]Info, error) {
	nodes, err := w.KubernetesLister.NodeLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var nodeInfos []Info
	for _, n := range nodes {
		nodeInfos = append(nodeInfos, Info{
			Name:              n.Name,
			Status:            checkNodeStatus(n.Status.Conditions),
			Roles:             checkNodeRole(n.Labels),
			Address:           getNodeAddress(n.Status.Addresses),
			OSImage:           n.Status.NodeInfo.OSImage,
			CreationTimestamp: n.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
		})
	}

	sort.Slice(nodeInfos, func(left, right int) bool {
		return nodeInfos[left].Name > nodeInfos[right].Name
	})

	return nodeInfos, nil
}

func checkNodeStatus(conditions []v1.NodeCondition) string {
	for _, condition := range conditions {
		if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
			return "Ready"
		}
	}

	return "NotReady"
}

func getNodeAddress(addresses []v1.NodeAddress) string {
	addressList := make([]string, 0)
	for _, addr := range addresses {
		addressList = append(addressList, addr.Address)
	}
	return strings.Join(addressList, "/")
}

func checkNodeRole(labels map[string]string) string {
	if _, ok := labels[masterRoleSpecialLabel]; ok {
		return "master"
	}

	return "worker"
}
