package kScrapper

import (
	"github.com/Masterlu1998/kube-viewer/kScrapper/resource"
	"github.com/Masterlu1998/kube-viewer/kube"
)

type ScrapperManagement struct {
	ResourceScrapper *resource.ResourceScrapper
}

func NewScrapperManagement() (*ScrapperManagement, error) {
	kubeClient, err := kube.GetKubernetesClient()
	if err != nil {
		return nil, err
	}

	s := &ScrapperManagement{
		ResourceScrapper: resource.NewResourceScrapper(kubeClient),
	}

	return s, nil
}
