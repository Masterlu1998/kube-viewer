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
	DeploymentScrapperTypes = "DeploymentScrapper"
	DeploymentResourceTypes = "Deployment"
)

type DeploymentScrapper struct {
	*common.CommonScrapper
}

func NewDeploymentScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *DeploymentScrapper {
	return &DeploymentScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *DeploymentScrapper) GetScrapperTypes() string {
	return DeploymentScrapperTypes
}

func (w *DeploymentScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)
}

func (w *DeploymentScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	deployments, err := getWorkloads(w.KubernetesClient, w.KubernetesLister, DeploymentResourceTypes, listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return deployments, nil
}
