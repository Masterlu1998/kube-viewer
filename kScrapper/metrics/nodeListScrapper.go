package metrics

import (
	"context"
	"sort"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

const (
	NodeMetricsListScrapperTypes = "NodeMetricsListScrapper"
)

// Node detail scrapper
type NodeMetricsListScrapper struct {
	*common.CommonScrapper
	kubernetesDynamicClient dynamic.Interface
}

func NewNodeMetricsListScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dynamicClient dynamic.Interface, dc *debug.DebugCollector) *NodeMetricsListScrapper {
	return &NodeMetricsListScrapper{
		kubernetesDynamicClient: dynamicClient,
		CommonScrapper:          common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *NodeMetricsListScrapper) GetScrapperTypes() string {
	return NodeMetricsListScrapperTypes
}

func (w *NodeMetricsListScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)
}

func (w *NodeMetricsListScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	nodeCurMetrics, err := w.kubernetesDynamicClient.Resource(nodeMetricsGVP).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	nodeMetricsMap := make(map[string]*NodeMetricsInfo)

	for _, m := range nodeCurMetrics.Items {
		metricItem, err := convertItemToNodeMetricInfo(m)
		if err != nil {
			continue
		}
		nodeMetricsMap[metricItem.Name] = metricItem
	}

	nodes, err := w.KubernetesLister.NodeLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	nodeOrderIndex := make([]string, 0)
	for _, n := range nodes {
		if metricItem, ok := nodeMetricsMap[n.Name]; ok {
			nodeOrderIndex = append(nodeOrderIndex, n.Name)
			memoryTotal := n.Status.Capacity.Memory()
			metricItem.MemoryTotal = memoryTotal
			cpuTotal := n.Status.Capacity.Cpu()
			metricItem.CPUTotal = cpuTotal
		}
	}

	sort.Strings(nodeOrderIndex)

	result := make([]*NodeMetricsInfo, 0)
	for _, nodeName := range nodeOrderIndex {
		result = append(result, nodeMetricsMap[nodeName])
	}

	return result, nil
}

func convertItemToNodeMetricInfo(m unstructured.Unstructured) (*NodeMetricsInfo, error) {
	cpuUsageStr, _, _ := unstructured.NestedString(m.Object, "usage", "cpu")
	cpuUsage, err := resource.ParseQuantity(cpuUsageStr)
	if err != nil {
		return nil, err
	}

	memoryUsageStr, _, _ := unstructured.NestedString(m.Object, "usage", "memory")
	memoryUsage, err := resource.ParseQuantity(memoryUsageStr)
	if err != nil {
		return nil, err
	}

	return &NodeMetricsInfo{
		Name:        m.GetName(),
		CPUUsage:    &cpuUsage,
		MemoryUsage: &memoryUsage,
	}, nil
}
