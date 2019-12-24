package pv

import (
	"sort"
	"time"

	"github.com/Masterlu1998/kube-viewer/kube"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	Name          string
	Capacity      string
	AccessMode    string
	ReclaimPolicy string
	Status        string
	StorageClass  string
	CreateTime    string
}

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
	kubernetesLister *kube.KubeLister
}

func (ka *kubeAccessor) getPVs(namespace string) ([]Info, error) {
	pvs, err := ka.kubernetesLister.PVLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var pvInfos []Info
	for _, p := range pvs {
		allAccessMode := ""
		for _, am := range p.Spec.AccessModes {
			allAccessMode += string(am) + " "
		}

		c := p.Spec.Capacity[v1.ResourceStorage]

		pvInfos = append(pvInfos, Info{
			Name:          p.Name,
			Capacity:      c.String(),
			AccessMode:    allAccessMode,
			ReclaimPolicy: string(p.Spec.PersistentVolumeReclaimPolicy),
			Status:        string(p.Status.Phase),
			StorageClass:  p.Spec.StorageClassName,
			CreateTime:    p.CreationTimestamp.Format(time.RFC3339),
		})
	}

	sort.Slice(pvInfos, func(left, right int) bool {
		return pvInfos[left].Name > pvInfos[right].Name
	})

	return pvInfos, nil
}
