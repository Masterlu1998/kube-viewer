package workload

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	JobDetailScrapperTypes = "JobDetailScrapper"
)

// Node detail scrapper
type JobDetailScrapper struct {
	*common.CommonScrapper
}

func NewJobDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *JobDetailScrapper {
	return &JobDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *JobDetailScrapper) GetScrapperTypes() string {
	return JobDetailScrapperTypes
}

func (w *JobDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *JobDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	job, err := w.CommonScrapper.KubernetesLister.JobLister.Jobs(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(job)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
