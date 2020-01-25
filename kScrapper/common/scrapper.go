package common

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const commonScrapperTypes = "commonScrapper"

type DataSourceFunc func(args ScrapperArgs) (KubernetesData, error)

func NewCommonScrapper(dc *debug.DebugCollector, client kubernetes.Interface, lister *kube.KubeLister) *CommonScrapper {
	return &CommonScrapper{
		namespace:      "",
		debugCollector: dc,
		// TODO: useless, need to remove
		KubernetesClient: client,
		KubernetesLister: lister,
	}
}

type CommonScrapper struct {
	stop             chan bool
	ch               chan KubernetesData
	namespace        string
	debugCollector   *debug.DebugCollector
	KubernetesClient kubernetes.Interface
	KubernetesLister *kube.KubeLister
}

func (c *CommonScrapper) Watch() <-chan KubernetesData {
	return c.ch
}

func (c *CommonScrapper) SetNamespace(ns string) {
	c.debugCollector.Collect(debug.NewDebugMessage(debug.Info,
		fmt.Sprintf("set namespace to %s", ns), "commonScrapper"))
	c.namespace = ns
}

func (c *CommonScrapper) ScrapeDataIntoChWithSource(ctx context.Context, f DataSourceFunc, args ScrapperArgs) {
	c.SetNamespace(args.GetNamespaceField())
	c.initScrapper()

	go func(ctx context.Context, stop chan bool, f DataSourceFunc) {
		ticker := time.NewTicker(ScrapInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				data, err := f(args)
				if err != nil {
					c.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), commonScrapperTypes))
					continue
				}
				c.ch <- data
			}
		}
	}(ctx, c.stop, f)
}

func (c *CommonScrapper) initScrapper() {
	c.StopScrapper()
	c.ch = make(chan KubernetesData)
	c.stop = make(chan bool)
}

func (c *CommonScrapper) StopScrapper() {
	if c.stop != nil {
		c.stop <- true
		close(c.stop)
	}
	c.stop = nil

	if c.ch != nil {
		close(c.ch)
	}
	c.ch = nil
}
