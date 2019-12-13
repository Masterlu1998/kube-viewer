package workload

import (
	"context"
	"time"

	"github.com/Masterlu1998/kube-viewer/dataTypes"
	"k8s.io/client-go/kubernetes"
)

const (
	ResourceScrapperTypes          = "DeploymentScrapper"
	DeploymentResourceTypes string = "Deployment"
)
const (
	StatefulSetResourceTypes = "StatefulSet"
	DaemonSetResourceTypes   = "DaemonSet"
	ReplicaSetResourceTypes  = "ReplicaSet"
	CronJobResourceTypes     = "CronJob"
	JobResourceTypes         = "Job"
)

type DeploymentScrapper struct {
	stop          chan bool
	ch            chan dataTypes.KubernetesData
	kubeAccessor  *kubeAccessor
	resourceTypes string
}

func NewDeploymentScrapper(client *kubernetes.Clientset) *DeploymentScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
	}

	return &DeploymentScrapper{
		kubeAccessor: ka,
	}
}

func (w *DeploymentScrapper) GetScrapperTypes() string {
	return ResourceScrapperTypes
}

func (w *DeploymentScrapper) Watch() <-chan dataTypes.KubernetesData {
	return w.ch
}

func (w *DeploymentScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.stopResourceScrapper()
	w.ch = make(chan dataTypes.KubernetesData)
	w.stop = make(chan bool)

	go func(ctx context.Context, stop chan bool) {
		ticker := time.NewTicker(dataTypes.ScrapInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = w.scrapeDataIntoCh(namespace)
			}
		}
	}(ctx, w.stop)
}

func (w *DeploymentScrapper) scrapeDataIntoCh(namespace string) error {
	deployments, err := w.kubeAccessor.getWorkloads(DeploymentResourceTypes, namespace)
	if err != nil {
		return err
	}

	w.ch <- deployments
	return nil
}

func (w *DeploymentScrapper) stopResourceScrapper() {
	if w.stop != nil {
		w.stop <- true
	}

	if w.ch != nil {
		close(w.ch)
	}
	w.ch = nil
}
