package kScrapper

import (
	"context"
	"errors"
	"sync"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	"github.com/Masterlu1998/kube-viewer/kube"
)

type Scrapper interface {
	GetScrapperTypes() string
	Watch() <-chan common.KubernetesData
	StartScrapper(ctx context.Context, namespace string)
	SetNamespace(namespace string)
	StopScrapper()
}

var workloadScrapperTypes = []string{
	workload.DeploymentScrapperTypes,
	workload.StatefulSetScrapperTypes,
	workload.DaemonSetScrapperTypes,
	workload.ReplicaSetScrapperTypes,
	workload.CronJobScrapperTypes,
	workload.JobScrapperTypes,
}

func NewScrapperManagement(ctx context.Context, collector *debug.DebugCollector) (*ScrapperManagement, error) {
	kubeClient, err := kube.GetKubernetesClient()
	if err != nil {
		return nil, err
	}

	kubeLister := kube.NewKubeLister(ctx, kubeClient)

	sMap := map[string]Scrapper{
		workload.DeploymentScrapperTypes:  workload.NewDeploymentScrapper(kubeLister, kubeClient, collector),
		workload.StatefulSetScrapperTypes: workload.NewStatefulSetScrapper(kubeLister, kubeClient, collector),
		workload.DaemonSetScrapperTypes:   workload.NewDaemonSetScrapper(kubeLister, kubeClient, collector),
		workload.ReplicaSetScrapperTypes:  workload.NewReplicaSetScrapper(kubeLister, kubeClient, collector),
		workload.CronJobScrapperTypes:     workload.NewCronJobScrapper(kubeLister, kubeClient, collector),
		workload.JobScrapperTypes:         workload.NewJobScrapper(kubeLister, kubeClient, collector),
		namespace.NamespaceScrapperTypes:  namespace.NewNamespaceScrapper(kubeLister, kubeClient, collector),
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

func (sm *ScrapperManagement) StartSpecificScrapper(ctx context.Context, scrapperType, namespace string) error {
	sm.rwMutex.Lock()
	defer sm.rwMutex.Unlock()
	if s, ok := sm.scrapperMap[scrapperType]; ok {
		s.StartScrapper(ctx, namespace)
		sm.activeScrapperMap[scrapperType] = true
		return nil
	}

	return errors.New("no this scrapper type")
}

func (sm *ScrapperManagement) StopSpecificScrapper(scrapperType string) {
	sm.rwMutex.Lock()
	defer sm.rwMutex.Unlock()
	if !sm.activeScrapperMap[scrapperType] {
		return
	}

	sm.scrapperMap[scrapperType].StopScrapper()
	delete(sm.activeScrapperMap, scrapperType)
}

func (sm *ScrapperManagement) StopAllWorkloadScrapper() {
	sm.rwMutex.Lock()
	defer sm.rwMutex.Unlock()

	for _, workloadScrapper := range workloadScrapperTypes {
		if sm.activeScrapperMap[workloadScrapper] {
			sm.scrapperMap[workloadScrapper].StopScrapper()
			delete(sm.activeScrapperMap, workloadScrapper)
		}
	}
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
