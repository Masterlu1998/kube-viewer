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
	DaemonSetScrapperTypes = "DaemonSetScrapper"
	DaemonSetResourceTypes = "DaemonSet"
)

type DaemonSetScrapper struct {
	*common.CommonScrapper
}

func NewDaemonSetScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *DaemonSetScrapper {
	return &DaemonSetScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (d *DaemonSetScrapper) GetScrapperTypes() string {
	return DaemonSetScrapperTypes
}

func (d *DaemonSetScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	d.CommonScrapper.ScrapeDataIntoChWithSource(ctx, d.scrapeDataIntoCh, args)
}

func (d *DaemonSetScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	daemonSets, err := getWorkloads(d.KubernetesClient, d.KubernetesLister, DaemonSetResourceTypes, listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return daemonSets, nil
}
