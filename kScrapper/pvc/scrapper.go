package pvc

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kube"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

const (
	PVCScrapperTypes = "PVCScrapper"
	PVCResourceTypes = "PVC"
)

type PVCScrapper struct {
	*common.CommonScrapper
}

func NewPVCScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *PVCScrapper {
	return &PVCScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *PVCScrapper) GetScrapperTypes() string {
	return PVCScrapperTypes
}

func (w *PVCScrapper) StartScrapper(ctx context.Context, args common.ScrapperArgs) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, args)

}

func (w *PVCScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	listArgs, ok := args.(common.ListScrapperArgs)
	if !ok {
		return nil, errors.New("convert to common.ListScrapperArgs failed")
	}

	configMapInfos, err := w.getPVCs(listArgs.Namespace)
	if err != nil {
		return nil, err
	}

	return configMapInfos, nil
}

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

func (w *PVCScrapper) getPVCs(namespace string) ([]Info, error) {
	pvcs, err := w.KubernetesLister.PVCLister.PersistentVolumeClaims(namespace).List(labels.Everything())
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
