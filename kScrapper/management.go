package kScrapper

import (
	"github.com/Masterlu1998/kube-viewer/dataTypes"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	"github.com/Masterlu1998/kube-viewer/kube"
)

type ScrapperManagement struct {
	ScrapperMap      map[string]dataTypes.Scrapper
	ResourceScrapper *workload.DeploymentScrapper
}

func NewScrapperManagement() (*ScrapperManagement, error) {
	kubeClient, err := kube.GetKubernetesClient()
	if err != nil {
		return nil, err
	}

	sMap := map[string]dataTypes.Scrapper{
		workload.ResourceScrapperTypes:   workload.NewDeploymentScrapper(kubeClient),
		namespace.NamespaceScrapperTypes: namespace.NewNamespaceScrapper(kubeClient),
	}

	return &ScrapperManagement{
		ScrapperMap: sMap,
	}, nil
}
