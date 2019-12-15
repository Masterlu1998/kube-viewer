package workload

import (
	"context"

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
	kubeAccessor *kubeAccessor
}

func NewDaemonSetScrapper(lister *kube.KubeLister, client *kubernetes.Clientset) *DaemonSetScrapper {
	return &DaemonSetScrapper{
		kubeAccessor:   generateKubeAccessor(lister, client),
		CommonScrapper: common.NewCommonScrapper(),
	}
}

func (d *DaemonSetScrapper) GetScrapperTypes() string {
	return DaemonSetScrapperTypes
}

func (d *DaemonSetScrapper) StartScrapper(ctx context.Context, namespace string) {
	d.CommonScrapper.ScrapeDataIntoChWithSource(ctx, d.scrapeDataIntoCh, namespace)
}

func (d *DaemonSetScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	daemonSets, err := d.kubeAccessor.getWorkloads(DaemonSetResourceTypes, namespace)
	if err != nil {
		return nil, err
	}

	return daemonSets, nil
}
