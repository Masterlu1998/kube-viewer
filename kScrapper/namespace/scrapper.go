package namespace

import (
	"context"
	"time"

	"github.com/Masterlu1998/kube-viewer/dataTypes"
	"k8s.io/client-go/kubernetes"
)

const (
	NamespaceScrapperTypes = "NamespaceScrapper"
	NamespaceResourceTypes = "namespace"
)

type NamespaceScrapper struct {
	stop          chan bool
	ch            chan dataTypes.KubernetesData
	kubeAccessor  *kubeAccessor
	resourceTypes string
}

func NewNamespaceScrapper(client *kubernetes.Clientset) *NamespaceScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
	}

	return &NamespaceScrapper{
		kubeAccessor: ka,
	}
}

func (w *NamespaceScrapper) GetScrapperTypes() string {
	return NamespaceScrapperTypes
}

func (w *NamespaceScrapper) Watch() <-chan dataTypes.KubernetesData {
	return w.ch
}

func (w *NamespaceScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.StopResourceScrapper()
	w.ch = make(chan dataTypes.KubernetesData)
	w.stop = make(chan bool)

	go func(ctx context.Context, stop chan bool) {
		ticker := time.NewTicker(dataTypes.ScrapInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = w.scrapeDataIntoCh()
			}
		}
	}(ctx, w.stop)
}

func (w *NamespaceScrapper) scrapeDataIntoCh() error {
	namespaceInfos, err := w.kubeAccessor.getNamespaces()
	if err != nil {
		return err
	}

	w.ch <- namespaceInfos
	return nil
}

func (w *NamespaceScrapper) StopResourceScrapper() {
	if w.stop != nil {
		w.stop <- true
	}

	if w.ch != nil {
		close(w.ch)
	}
	w.ch = nil
}
