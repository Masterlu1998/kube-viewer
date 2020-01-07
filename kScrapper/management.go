package kScrapper

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/configMap"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/node"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pv"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pvc"
	"github.com/Masterlu1998/kube-viewer/kScrapper/secret"
	"github.com/Masterlu1998/kube-viewer/kScrapper/service"
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

var mainScrapperTypes = []string{
	workload.DeploymentScrapperTypes,
	workload.StatefulSetScrapperTypes,
	workload.DaemonSetScrapperTypes,
	workload.ReplicaSetScrapperTypes,
	workload.CronJobScrapperTypes,
	workload.JobScrapperTypes,
	service.ServiceScrapperTypes,
	configMap.ConfigMapScrapperTypes,
	secret.SecretScrapperTypes,
	pvc.PVCScrapperTypes,
	pv.PVScrapperTypes,
	node.NodeScrapperTypes,
}

func NewScrapperManagement(ctx context.Context, collector *debug.DebugCollector) (*ScrapperManagement, error) {
	kubeClient, err := kube.GetKubernetesClient()
	if err != nil {
		return nil, err
	}

	kubeLister := kube.NewKubeLister(ctx, kubeClient)

	// TODO: I will make a factory to create scrapper
	sMap := map[string]Scrapper{
		workload.DeploymentScrapperTypes:  workload.NewDeploymentScrapper(kubeLister, kubeClient, collector),
		workload.StatefulSetScrapperTypes: workload.NewStatefulSetScrapper(kubeLister, kubeClient, collector),
		workload.DaemonSetScrapperTypes:   workload.NewDaemonSetScrapper(kubeLister, kubeClient, collector),
		workload.ReplicaSetScrapperTypes:  workload.NewReplicaSetScrapper(kubeLister, kubeClient, collector),
		workload.CronJobScrapperTypes:     workload.NewCronJobScrapper(kubeLister, kubeClient, collector),
		workload.JobScrapperTypes:         workload.NewJobScrapper(kubeLister, kubeClient, collector),
		service.ServiceScrapperTypes:      service.NewNamespaceScrapper(kubeLister, kubeClient, collector),
		configMap.ConfigMapScrapperTypes:  configMap.NewConfigMapScrapper(kubeLister, kubeClient, collector),
		secret.SecretScrapperTypes:        secret.NewSecretScrapper(kubeLister, kubeClient, collector),
		pvc.PVCScrapperTypes:              pvc.NewPVCScrapper(kubeLister, kubeClient, collector),
		pv.PVScrapperTypes:                pv.NewPVScrapper(kubeLister, kubeClient, collector),
		node.NodeScrapperTypes:            node.NewNodeScrapper(kubeLister, kubeClient, collector),

		namespace.NamespaceScrapperTypes: namespace.NewNamespaceScrapper(kubeLister, kubeClient, collector),
	}

	return &ScrapperManagement{
		scrapperMap:       sMap,
		activeScrapperMap: make(map[string]bool),
		namespace:         "",
		debugCollector:    collector,
	}, nil
}

type ScrapperManagement struct {
	rwMutex           sync.RWMutex
	scrapperMap       map[string]Scrapper
	activeScrapperMap map[string]bool
	namespace         string
	debugCollector    *debug.DebugCollector
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

// TODO: maybe useless
func (sm *ScrapperManagement) StopSpecificScrapper(scrapperType string) {
	sm.rwMutex.Lock()
	defer sm.rwMutex.Unlock()
	if !sm.activeScrapperMap[scrapperType] {
		return
	}

	sm.scrapperMap[scrapperType].StopScrapper()
	delete(sm.activeScrapperMap, scrapperType)
}

func (sm *ScrapperManagement) StopMainScrapper() {
	sm.rwMutex.Lock()
	defer sm.rwMutex.Unlock()

	for _, workloadScrapper := range mainScrapperTypes {
		if sm.activeScrapperMap[workloadScrapper] {
			sm.debugCollector.Collect(debug.NewDebugMessage(debug.Info, fmt.Sprintf("stop %s", workloadScrapper), "scrapperManagement"))
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
