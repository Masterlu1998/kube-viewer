package workload

import (
	"context"

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
	kubeAccessor *kubeAccessor
}

func NewDeploymentScrapper(lister *kube.KubeLister, client *kubernetes.Clientset) *DeploymentScrapper {
	return &DeploymentScrapper{
		kubeAccessor:   generateKubeAccessor(lister, client),
		CommonScrapper: common.NewCommonScrapper(),
	}
}

func (w *DeploymentScrapper) GetScrapperTypes() string {
	return DeploymentScrapperTypes
}

func (w *DeploymentScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, namespace)
}

func (w *DeploymentScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	deployments, err := w.kubeAccessor.getWorkloads(DeploymentResourceTypes, namespace)
	if err != nil {
		return nil, err
	}

	return deployments, nil
}
