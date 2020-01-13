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
	JobScrapperTypes = "JobScrapper"
	JobResourceTypes = "Job"
)

type JobScrapper struct {
	*common.CommonScrapper
}

func NewJobListScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *JobScrapper {
	return &JobScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (j *JobScrapper) GetScrapperTypes() string {
	return JobScrapperTypes
}

func (j *JobScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	j.CommonScrapper.ScrapeDataIntoChWithSource(ctx, j.scrapeDataIntoCh, args)
}

func (j *JobScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	jobs, err := getWorkloads(j.KubernetesClient, j.KubernetesLister, StatefulSetResourceTypes, listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
