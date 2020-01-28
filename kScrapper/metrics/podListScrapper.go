package metrics

import (
	"context"
	"errors"
	"sort"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

const (
	PodMetricsListScrapperTypes = "PodMetricsListScrapper"
)

// Pod detail scrapper
type PodMetricsListScrapper struct {
	*common.CommonScrapper
	kubernetesDynamicClient dynamic.Interface
}

func NewPodMetricsListScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dynamicClient dynamic.Interface, dc *debug.DebugCollector) *PodMetricsListScrapper {
	return &PodMetricsListScrapper{
		kubernetesDynamicClient: dynamicClient,
		CommonScrapper:          common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *PodMetricsListScrapper) GetScrapperTypes() string {
	return PodMetricsListScrapperTypes
}

func (w *PodMetricsListScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)
}

func (w *PodMetricsListScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	podMetrics, err := w.kubernetesDynamicClient.Resource(podMetricsGVP).Namespace(listArgs.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]*PodMetricsInfo, 0)
	for _, m := range podMetrics.Items {
		metricsInfo, err := convertItemToPodMetricInfo(m)
		if err != nil {
			return nil, err
		}

		result = append(result, metricsInfo)
	}

	sort.Slice(result, func(left, right int) bool {
		return result[left].CPUUsage.Cmp(*result[right].CPUUsage) == 1
	})

	return result, nil
}

func convertItemToPodMetricInfo(m unstructured.Unstructured) (*PodMetricsInfo, error) {
	containers, _, _ := unstructured.NestedSlice(m.Object, "containers")
	totalCPUUsage, err := resource.ParseQuantity("0")
	if err != nil {
		return nil, err
	}

	totalMemoryUsage, err := resource.ParseQuantity("0")
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		obj, ok := c.(map[string]interface{})
		if !ok {
			return nil, errors.New("parse to map[string]interface failed")
		}

		rawCPUPerContainer, _, _ := unstructured.NestedString(obj, "usage", "cpu")
		rawMemoryPerContainer, _, _ := unstructured.NestedString(obj, "usage", "memory")

		cpuPerContainer, err := resource.ParseQuantity(rawCPUPerContainer)
		if err != nil {
			return nil, err
		}

		memoryPerContainer, err := resource.ParseQuantity(rawMemoryPerContainer)
		if err != nil {
			return nil, err
		}

		totalCPUUsage.Add(cpuPerContainer)
		totalMemoryUsage.Add(memoryPerContainer)
	}

	metricsInfo := &PodMetricsInfo{
		Name:        m.GetName(),
		NameSpace:   m.GetNamespace(),
		CPUUsage:    &totalCPUUsage,
		MemoryUsage: &totalMemoryUsage,
	}
	return metricsInfo, nil
}
