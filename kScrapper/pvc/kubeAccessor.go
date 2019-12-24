package pvc

import (
	"sort"
	"time"

	"github.com/Masterlu1998/kube-viewer/kube"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	Name         string
	Namespace    string
	Status       string
	Volume       string
	Request      string
	Limit        string
	AccessMode   string
	StorageClass string
	CreateTime   string
}

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
	kubernetesLister *kube.KubeLister
}

func (ka *kubeAccessor) getPVCs(namespace string) ([]Info, error) {
	pvcs, err := ka.kubernetesLister.PVCLister.PersistentVolumeClaims(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var pvcInfos []Info
	for _, s := range pvcs {
		allAccessMode := ""
		for _, am := range s.Spec.AccessModes {
			allAccessMode += string(am) + " "
		}
		r, l := s.Spec.Resources.Requests[v1.ResourceStorage], s.Spec.Resources.Limits[v1.ResourceStorage]

		sClass := ""
		if s.Spec.StorageClassName != nil {
			sClass = *s.Spec.StorageClassName
		}

		pvcInfos = append(pvcInfos, Info{
			Name:         s.Name,
			Namespace:    s.Namespace,
			Status:       string(s.Status.Phase),
			Volume:       s.Spec.VolumeName,
			Request:      r.String(),
			Limit:        l.String(),
			AccessMode:   allAccessMode,
			StorageClass: sClass,
			CreateTime:   s.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
		})
	}

	sort.Slice(pvcInfos, func(left, right int) bool {
		return pvcInfos[left].Name > pvcInfos[right].Name
	})

	return pvcInfos, nil
}
