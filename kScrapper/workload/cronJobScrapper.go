package workload

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	CronJobScrapperTypes = "CronJobScrapper"
	CronJobResourceTypes = "CronJob"
)

type CronJobScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewCronJobScrapper(lister *kube.KubeLister, client *kubernetes.Clientset) *CronJobScrapper {
	return &CronJobScrapper{
		kubeAccessor:   generateKubeAccessor(lister, client),
		CommonScrapper: common.NewCommonScrapper(),
	}
}

func (c *CronJobScrapper) GetScrapperTypes() string {
	return CronJobScrapperTypes
}

func (c *CronJobScrapper) StartScrapper(ctx context.Context, namespace string) {
	c.CommonScrapper.ScrapeDataIntoChWithSource(ctx, c.scrapeDataIntoCh, namespace)
}

func (c *CronJobScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	cronJobs, err := c.kubeAccessor.getWorkloads(CronJobResourceTypes, namespace)
	if err != nil {
		return nil, err
	}

	return cronJobs, nil
}
