package kScrapper

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/dataTypes"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	"github.com/Masterlu1998/kube-viewer/kube"
)

func NewScrapperManagement() (*ScrapperManagement, error) {
	kubeClient, err := kube.GetKubernetesClient()
	if err != nil {
		return nil, err
	}

	sMap := map[string]dataTypes.Scrapper{
		workload.ResourceScrapperTypes:   workload.NewDeploymentScrapper(kubeClient, ""),
		namespace.NamespaceScrapperTypes: namespace.NewNamespaceScrapper(kubeClient, ""),
	}

	return &ScrapperManagement{
		ScrapperMap:       sMap,
		activeScrapperMap: make(map[string]bool),
		namespace:         "",
	}, nil
}

type ScrapperManagement struct {
	ScrapperMap       map[string]dataTypes.Scrapper
	activeScrapperMap map[string]bool
	namespace         string
}

func (sm *ScrapperManagement) StartSpecificScrapper(ctx context.Context, scrapperType string) {
	sm.ScrapperMap[scrapperType].StartScrapper(ctx)
	sm.activeScrapperMap[scrapperType] = true
}

func (sm *ScrapperManagement) GetSpecificScrapperCh(scrapperType string) <-chan dataTypes.KubernetesData {
	return sm.ScrapperMap[scrapperType].Watch()
}

func (sm *ScrapperManagement) ResetNamespace(ns string) {
	for st := range sm.activeScrapperMap {
		// no need for namespace scrapper
		if st == namespace.NamespaceScrapperTypes {
			continue
		}
		sm.ScrapperMap[st].SetNamespace(ns)
	}
}
