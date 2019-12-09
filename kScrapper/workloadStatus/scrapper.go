package workloadStatus

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

const WorkloadStatusTypes = "workloadStatus"

type chData struct {
	deployments []deploymentInfo
}

type deploymentInfo struct {
	Name       string
	Namespace  string
	PodsLive   int
	PodsTotal  int
	CreateTime string
	Images     string
}

type workloadStatusScrapper struct {
	ch           chan common.Data
	kubeAccessor *kubeAccessor
}

func NewWorkloadStatusScrapper(client *kubernetes.Clientset) *workloadStatusScrapper {
	ka := &kubeAccessor{
		kubernetesClient: client,
	}
	return &workloadStatusScrapper{
		ch:           make(chan common.Data),
		kubeAccessor: ka,
	}
}

func (w *workloadStatusScrapper) GetScrapperTypes() string {
	return WorkloadStatusTypes
}

func (w *workloadStatusScrapper) ScrapeDataIntoCh() error {
	data := chData{}

	deployments, err := w.kubeAccessor.getDeployments()
	if err != nil {
		return err
	}

	data.deployments = deployments

	w.ch <- data
	return nil
}

func (w *workloadStatusScrapper) GraphAction(l *widgets.List) error {
	for d := range w.ch {
		workloadSData, ok := d.(chData)
		if !ok {
			logrus.Error("invalid chData type \"chData\"")
			return errors.New("invalid chData type \"chData\"")
		}

		for _, wd := range workloadSData.deployments {
			l.Rows = append(l.Rows, wd.Name)
		}
		ui.Render(l)
	}

	return nil
}

func (w *workloadStatusScrapper) GetDataCh() <-chan common.Data {
	return w.ch
}
