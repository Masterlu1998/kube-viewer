package common

import (
	"context"
	"time"
)

type DataSourceFunc func(string) (KubernetesData, error)

func NewCommonScrapper() *CommonScrapper {
	return &CommonScrapper{
		namespace: "",
	}
}

type CommonScrapper struct {
	stop      chan bool
	ch        chan KubernetesData
	namespace string
}

func (c *CommonScrapper) initScrapper() {
	c.stopResourceScrapper()
	c.ch = make(chan KubernetesData)
	c.stop = make(chan bool)
}

func (c *CommonScrapper) stopResourceScrapper() {
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

func (c *CommonScrapper) Watch() <-chan KubernetesData {
	return c.ch
}

func (c *CommonScrapper) SetNamespace(ns string) {
	c.namespace = ns
}

func (c *CommonScrapper) ScrapeDataIntoChWithSource(ctx context.Context, f DataSourceFunc) {
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
				data, err := f(c.namespace)
				if err != nil {
					continue
				}
				c.ch <- data
			}
		}
	}(ctx, c.stop, f)
}
