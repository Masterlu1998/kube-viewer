package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

const (
	ServiceScrapperTypes = "ServiceScrapper"
	ServiceResourceTypes = "Namespace"
)

type ServiceScrapper struct {
	*common.CommonScrapper
}

func NewNServiceListScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *ServiceScrapper {
	return &ServiceScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *ServiceScrapper) GetScrapperTypes() string {
	return ServiceScrapperTypes
}

func (w *ServiceScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *ServiceScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	namespaceInfos, err := w.getServices(listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return namespaceInfos, nil
}

type Info struct {
	Name      string
	Namespace string
	ClusterIP string
	Port      string
}

func (w *ServiceScrapper) getServices(namespace string) ([]Info, error) {
	services, err := w.KubernetesLister.ServiceLister.Services(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var serviceInfos []Info
	for _, s := range services {
		ports := ""
		for _, p := range s.Spec.Ports {
			var itemPort string
			if p.NodePort != 0 {
				itemPort = fmt.Sprintf("%d->%d->%d/%s ", p.NodePort, p.Port, p.TargetPort.IntValue(), p.Protocol)
			} else {
				itemPort = fmt.Sprintf("%d->%d/%s ", p.Port, p.TargetPort.IntValue(), p.Protocol)
			}
			ports += itemPort
		}

		serviceInfos = append(serviceInfos, Info{
			Name:      s.Name,
			Namespace: s.Namespace,
			ClusterIP: s.Spec.ClusterIP,
			Port:      ports,
		})
	}

	sort.Slice(serviceInfos, func(left, right int) bool {
		return serviceInfos[left].Name > serviceInfos[right].Name
	})

	return serviceInfos, nil
}
