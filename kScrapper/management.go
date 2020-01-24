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
	StartScrapper(ctx context.Context, args common.ScrapperArgs)
	SetNamespace(namespace string)
	StopScrapper()
}

func NewScrapperManagement(ctx context.Context, collector *debug.DebugCollector) (*ScrapperManagement, error) {
	kubeClient, err := kube.GetKubernetesClient()
	if err != nil {
		return nil, err
	}

	kubeLister := kube.NewKubeLister(ctx, kubeClient)

	// TODO: I will make a factory to create scrapper
	sMap := map[string]Scrapper{
		workload.DeploymentListScrapperTypes:  workload.NewDeploymentListScrapper(kubeLister, kubeClient, collector),
		workload.StatefulSetListScrapperTypes: workload.NewStatefulSetListScrapper(kubeLister, kubeClient, collector),
		workload.DaemonSetListScrapperTypes:   workload.NewDaemonSetListScrapper(kubeLister, kubeClient, collector),
		workload.ReplicaSetListScrapperTypes:  workload.NewReplicaSetListScrapper(kubeLister, kubeClient, collector),
		workload.CronJobListScrapperTypes:     workload.NewCronJobListScrapper(kubeLister, kubeClient, collector),
		workload.JobListScrapperTypes:         workload.NewJobListScrapper(kubeLister, kubeClient, collector),
		service.ServiceListScrapperTypes:      service.NewNServiceListScrapper(kubeLister, kubeClient, collector),
		configMap.ConfigListMapScrapperTypes:  configMap.NewConfigMapListScrapper(kubeLister, kubeClient, collector),
		secret.SecretListScrapperTypes:        secret.NewSecretListScrapper(kubeLister, kubeClient, collector),
		pvc.PVCListScrapperTypes:              pvc.NewPVCListScrapper(kubeLister, kubeClient, collector),
		pv.PVListScrapperTypes:                pv.NewPVListScrapper(kubeLister, kubeClient, collector),
		node.NodeListScrapperTypes:            node.NewNodeListScrapper(kubeLister, kubeClient, collector),

		workload.DeploymentDetailScrapperTypes:  workload.NewDeploymentDetailScrapper(kubeLister, kubeClient, collector),
		workload.StatefulSetDetailScrapperTypes: workload.NewStatefulSetDetailScrapper(kubeLister, kubeClient, collector),
		workload.DaemonSetDetailScrapperTypes:   workload.NewDaemonSetDetailScrapper(kubeLister, kubeClient, collector),
		workload.ReplicaSetDetailScrapperTypes:  workload.NewReplicaSetDetailScrapper(kubeLister, kubeClient, collector),
		workload.CronJobDetailScrapperTypes:     workload.NewCronJobDetailScrapper(kubeLister, kubeClient, collector),
		workload.JobDetailScrapperTypes:         workload.NewJobDetailScrapper(kubeLister, kubeClient, collector),
		service.ServiceDetailScrapperTypes:      service.NewServiceDetailScrapper(kubeLister, kubeClient, collector),
		configMap.ConfigMapDetailScrapperTypes:  configMap.NewConfigMapDetailScrapper(kubeLister, kubeClient, collector),
		secret.SecretDetailScrapperTypes:        secret.NewSecretDetailScrapper(kubeLister, kubeClient, collector),
		pvc.PVCDetailScrapperTypes:              pvc.NewPVCDetailScrapper(kubeLister, kubeClient, collector),
		pv.PVDetailScrapperTypes:                pv.NewPVDetailScrapper(kubeLister, kubeClient, collector),
		node.NodeDetailScrapperTypes:            node.NewNodeDetailScrapper(kubeLister, kubeClient, collector),

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
	mutex             sync.Mutex
	scrapperMap       map[string]Scrapper
	activeScrapperMap map[string]bool
	namespace         string
	debugCollector    *debug.DebugCollector
}

func (sm *ScrapperManagement) StartSpecificScrapper(ctx context.Context, scrapperType string, args common.ScrapperArgs) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if s, ok := sm.scrapperMap[scrapperType]; ok {
		s.StartScrapper(ctx, args)
		sm.activeScrapperMap[scrapperType] = true
		return nil
	}

	return errors.New("no this scrapper type")
}

func (sm *ScrapperManagement) StopMainScrapper() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for scrapper := range sm.activeScrapperMap {
		if scrapper == namespace.NamespaceScrapperTypes {
			continue
		}

		sm.debugCollector.Collect(debug.NewDebugMessage(debug.Info, fmt.Sprintf("stop %s", scrapper), "ScrapperManagement"))
		sm.scrapperMap[scrapper].StopScrapper()
		delete(sm.activeScrapperMap, scrapper)
	}
}

func (sm *ScrapperManagement) GetSpecificScrapperCh(scrapperType string) <-chan common.KubernetesData {
	return sm.scrapperMap[scrapperType].Watch()
}

func (sm *ScrapperManagement) ResetNamespace(ns string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	for st := range sm.activeScrapperMap {
		// no need for namespace scrapper
		if st == namespace.NamespaceScrapperTypes {
			continue
		}
		sm.scrapperMap[st].SetNamespace(ns)
	}
}
