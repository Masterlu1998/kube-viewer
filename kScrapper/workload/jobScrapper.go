package workload

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
)

const (
	JobScrapperTypes = "JobScrapper"
	JobResourceTypes = "Job"
)

type JobScrapper struct {
	*common.CommonScrapper
	kubeAccessor *kubeAccessor
}

func NewJobScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *JobScrapper {
	return &JobScrapper{
		kubeAccessor:   generateKubeAccessor(lister, client),
		CommonScrapper: common.NewCommonScrapper(dc),
	}
}

func (c *JobScrapper) GetScrapperTypes() string {
	return JobScrapperTypes
}

func (c *JobScrapper) StartScrapper(ctx context.Context, namespace string) {
	c.CommonScrapper.ScrapeDataIntoChWithSource(ctx, c.scrapeDataIntoCh, namespace)
}

func (c *JobScrapper) scrapeDataIntoCh(namespace string) (common.KubernetesData, error) {
	jobs, err := c.kubeAccessor.getWorkloads(JobResourceTypes, namespace)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
