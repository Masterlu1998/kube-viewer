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
	CronJobDetailScrapperTypes = "CronJobDetailScrapper"
)

// Node detail scrapper
type CronJobDetailScrapper struct {
	*common.CommonScrapper
}

func NewCronJobDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *CronJobDetailScrapper {
	return &CronJobDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *CronJobDetailScrapper) GetScrapperTypes() string {
	return CronJobDetailScrapperTypes
}

func (w *CronJobDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *CronJobDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	cronJob, err := w.CommonScrapper.KubernetesLister.CronJobLister.CronJobs(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(cronJob)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
