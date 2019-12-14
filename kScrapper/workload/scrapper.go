package workload

import (
	"context"
	"time"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
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
	ch            chan common.KubernetesData
	kubeAccessor  *kubeAccessor
	resourceTypes string
	namespace     string
}

func NewDeploymentScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, namespace string) *DeploymentScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
		kubernetesLister: lister,
	}

	return &DeploymentScrapper{
		kubeAccessor: ka,
		namespace:    namespace,
	}
}

func (w *DeploymentScrapper) GetScrapperTypes() string {
	return ResourceScrapperTypes
}

func (w *DeploymentScrapper) Watch() <-chan common.KubernetesData {
	return w.ch
}

func (w *DeploymentScrapper) StartScrapper(ctx context.Context) {
	w.stopResourceScrapper()
	w.ch = make(chan common.KubernetesData)
	w.stop = make(chan bool)

	go func(ctx context.Context, stop chan bool) {
		ticker := time.NewTicker(common.ScrapInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = w.scrapeDataIntoCh(w.namespace)
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

func (w *DeploymentScrapper) SetNamespace(namespace string) {
	w.namespace = namespace
}
