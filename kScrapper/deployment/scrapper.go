package deployment

import (
	"github.com/Masterlu1998/kube-viewer/dataTypes"
	"k8s.io/client-go/kubernetes"
)

type deploymentScrapper struct {
	ch           chan dataTypes.Data
	kubeAccessor *kubeAccessor
}

func NewDeploymentScrapper(client *kubernetes.Clientset) *deploymentScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
	}
	return &deploymentScrapper{
		ch:           make(chan dataTypes.Data),
		kubeAccessor: ka,
	}
}

func (w *deploymentScrapper) GetScrapperTypes() string {
	return dataTypes.DeploymentScrapperTypes
}

func (w *deploymentScrapper) ScrapeDataIntoCh() error {
	data := dataTypes.DeploymentScrapperChData{}

	deployments, err := w.kubeAccessor.getDeployments()
	if err != nil {
		return err
	}

	data.Deployments = deployments

	w.ch <- data
	return nil
}

func (w *deploymentScrapper) GetDataCh() <-chan dataTypes.Data {
	return w.ch
}
