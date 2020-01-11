package pv

import (
	"context"
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
	PVScrapperTypes = "PVScrapper"
	PVResourceTypes = "PV"
)

type PVScrapper struct {
	*common.CommonScrapper
}

func NewPVScrapper(lister *kube.KubeLister, client *kubernetes.Clientset, dc *debug.DebugCollector) *PVScrapper {
	return &PVScrapper{
		CommonScrapper: common.NewCommonScrapper(dc, client, lister),
	}
}

func (w *PVScrapper) GetScrapperTypes() string {
	return PVScrapperTypes
}

func (w *PVScrapper) StartScrapper(ctx context.Context, namespace string) {
	w.CommonScrapper.ScrapeDataIntoChWithSource(ctx, w.scrapeDataIntoCh, namespace)

}

func (w *PVScrapper) scrapeDataIntoCh(args common.ScrapperArgs) (common.KubernetesData, error) {
	PVInfos, err := w.getPVs()
	if err != nil {
		return nil, err
	}

	return PVInfos, nil
}

type Info struct {
	Name          string
	Capacity      string
	AccessMode    string
	ReclaimPolicy string
	Status        string
	StorageClass  string
	CreateTime    string
}

func (w *PVScrapper) getPVs() ([]Info, error) {
	pvs, err := w.KubernetesLister.PVLister.List(labels.Everything())
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
