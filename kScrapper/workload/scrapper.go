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
const (
	StatefulSetResourceTypes = "StatefulSet"
	DaemonSetResourceTypes   = "DaemonSet"
	ReplicaSetResourceTypes  = "ReplicaSet"
	CronJobResourceTypes     = "CronJob"
	JobResourceTypes         = "Job"
)

type DeploymentScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewDeploymentScrapper(lister *kube.KubeLister, client *kubernetes.Clientset) *DeploymentScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
		kubernetesLister: lister,
	}

	return &DeploymentScrapper{
		kubeAccessor:   ka,
		CommonScrapper: common.NewCommonScrapper(),
	}
}

func (w *DeploymentScrapper) GetScrapperTypes() string {
	return DeploymentResourceTypes
}

func (w *DeploymentScrapper) StartScrapper(ctx context.Context) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh)
}

func (w *DeploymentScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	deployments, err := w.kubeAccessor.getWorkloads(DeploymentResourceTypes, namespace)
	if err != nil {
		return nil, err
	}

	return deployments, nil
}
