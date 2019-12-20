package service

import (
	"fmt"
	"sort"

	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	Name      string
	Namespace string
	ClusterIP string
	Port      string
}

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
	kubernetesLister *kube.KubeLister
}

func (ka *kubeAccessor) getServices(namespace string) ([]Info, error) {
	services, err := ka.kubernetesLister.ServiceLister.Services(namespace).List(labels.Everything())
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
