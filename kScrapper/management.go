package kScrapper

import (
	"context"
	"sync"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	"github.com/Masterlu1998/kube-viewer/kube"
)

type Scrapper interface {
	GetScrapperTypes() string
	Watch() <-chan common.KubernetesData
	StartScrapper(ctx context.Context)
	SetNamespace(namespace string)
}

func NewScrapperManagement() (*ScrapperManagement, error) {
	kubeClient, err := kube.GetKubernetesClient()
	if err != nil {
		return nil, err
	}

	sMap := map[string]Scrapper{
		workload.ResourceScrapperTypes:   workload.NewDeploymentScrapper(kubeClient, ""),
		namespace.NamespaceScrapperTypes: namespace.NewNamespaceScrapper(kubeClient, ""),
	}

	return &ScrapperManagement{
		scrapperMap:       sMap,
		activeScrapperMap: make(map[string]bool),
		namespace:         "",
	}, nil
}

type ScrapperManagement struct {
	rwMutex           sync.RWMutex
	scrapperMap       map[string]Scrapper
	activeScrapperMap map[string]bool
	namespace         string
}

func (sm *ScrapperManagement) StartSpecificScrapper(ctx context.Context, scrapperType string) {
	sm.rwMutex.Lock()
	defer sm.rwMutex.Unlock()
	sm.scrapperMap[scrapperType].StartScrapper(ctx)
	sm.activeScrapperMap[scrapperType] = true
}

func (sm *ScrapperManagement) GetSpecificScrapperCh(scrapperType string) <-chan common.KubernetesData {
	return sm.scrapperMap[scrapperType].Watch()
}

func (sm *ScrapperManagement) ResetNamespace(ns string) {
	sm.rwMutex.RLock()
	defer sm.rwMutex.RUnlock()
	for st := range sm.activeScrapperMap {
		// no need for namespace scrapper
		if st == namespace.NamespaceScrapperTypes {
			continue
		}
		sm.scrapperMap[st].SetNamespace(ns)
	}
}
