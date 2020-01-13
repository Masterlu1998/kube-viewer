package workload

import (
	"context"
	"errors"

	"github.com/Masterlu1998/kube-viewer/debug"
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
}

func NewCronJobScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *CronJobScrapper {
	return &CronJobScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (c *CronJobScrapper) GetScrapperTypes() string {
	return CronJobScrapperTypes
}

func (c *CronJobScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	c.CommonScrapper.ScrapeDataIntoChWithSource(ctx, c.scrapeDataIntoCh, args)
}

func (c *CronJobScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	cronJobs, err := getWorkloads(c.KubernetesClient, c.KubernetesLister, CronJobResourceTypes, listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return cronJobs, nil
}
