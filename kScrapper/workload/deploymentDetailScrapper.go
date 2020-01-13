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
	DeploymentDetailScrapperTypes = "DeploymentDetailScrapper"
)

// Node detail scrapper
type DeploymentDetailScrapper struct {
	*common.CommonScrapper
}

func NewDeploymentDetailScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *DeploymentDetailScrapper {
	return &DeploymentDetailScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *DeploymentDetailScrapper) GetScrapperTypes() string {
	return DeploymentDetailScrapperTypes
}

func (w *DeploymentDetailScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *DeploymentDetailScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	detailArgs := args.(common.DetailScrapperArgs)
	deployment, err := w.CommonScrapper.KubernetesLister.DeploymentLister.Deployments(detailArgs.Namespace).Get(detailArgs.Name)
	if err != nil {
		return "", err
	}

	nodeYaml, err := yaml.Marshal(deployment)
	if err != nil {
		return "", err
	}

	return string(nodeYaml), nil
}
