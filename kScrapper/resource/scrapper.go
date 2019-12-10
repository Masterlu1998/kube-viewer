package resource

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
)

type kubernetesData interface{}

const (
	ResourceScrapperTypes = "resourceScrapper"
	scrapInterval         = time.Second * 1
)

type ResourceTypes string

const (
	DeploymentResourceTypes  ResourceTypes = "deployment"
	StatefulSetResourceTypes ResourceTypes = "statefulSet"
	DaemonSetResourceTypes   ResourceTypes = "daemonSet"
	ReplicaSetResourceTypes  ResourceTypes = "replicaSet"
	CronJobResourceTypes     ResourceTypes = "cronJob"
	JobResourceTypes         ResourceTypes = "jobResource"
)

type WorkloadData struct {
	Infos []WorkloadInfo
}

type WorkloadInfo struct {
	Kind       string
	Name       string
	Namespace  string
	PodsLive   string
	PodsTotal  string
	CreateTime string
	Images     string
}

type ResourceScrapper struct {
	ownCtx        context.Context
	ownCancel     context.CancelFunc
	ch            chan kubernetesData
	kubeAccessor  *kubeAccessor
	resourceTypes ResourceTypes
}

func NewResourceScrapper(client *kubernetes.Clientset) *ResourceScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
	}

	return &ResourceScrapper{
		kubeAccessor: ka,
	}
}

func (w *ResourceScrapper) GetScrapperTypes() string {
	return ResourceScrapperTypes
}

func (w *ResourceScrapper) GetDataCh() <-chan kubernetesData {
	return w.ch
}

func (w *ResourceScrapper) StartResourceScrapper(ctx context.Context, types ResourceTypes) {
	w.ch = make(chan kubernetesData)
	ownCtx, ownCancel := context.WithCancel(ctx)
	w.ownCtx, w.ownCancel = ownCtx, ownCancel
	go func() {
		ticker := time.NewTicker(scrapInterval)
		defer ticker.Stop()
		for {
			select {
			case <-w.ownCtx.Done():
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = w.scrapeDataIntoCh(types)
			}
		}
	}()
}

func (w *ResourceScrapper) scrapeDataIntoCh(resourceTypes ResourceTypes) error {
	workloads, err := w.kubeAccessor.getWorkloads(resourceTypes)
	if err != nil {
		return err
	}

	data := WorkloadData{Infos: workloads}

	w.ch <- data
	return nil
}

func (w *ResourceScrapper) StopResourceScrapper() {
	if w.ownCancel != nil {
		w.ownCancel()
	}
	w.ownCtx, w.ownCancel = nil, nil

	if w.ch != nil {
		close(w.ch)
	}
	w.ch = nil
}
