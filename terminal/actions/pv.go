package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pv"
)

var (
	pvTableHeader   = []string{"name", "capacity", "accessMode", "reclaim policy", "status", "storage class", "create time"}
	pvTableColWidth = []int{15, 10, 20, 20, 15, 20, 30}
)

func pvListDataGetter(c common.KubernetesData) ([]string, [][]string, []int, error) {
	pvInfos, ok := c.([]pv.Info)
	if !ok {
		return nil, nil, nil, errors.New("convert to pv.Info failed")
	}

	newPVTableData := make([][]string, 0)
	for _, pvInfo := range pvInfos {
		newPVTableData = append(newPVTableData, []string{
			pvInfo.Name,
			pvInfo.Capacity,
			pvInfo.AccessMode,
			pvInfo.ReclaimPolicy,
			pvInfo.Status,
			pvInfo.StorageClass,
			pvInfo.CreateTime,
		})
	}
	return pvTableHeader, newPVTableData, pvTableColWidth, nil
}

func BuildPVListAction() ActionHandler {
	return listResourceAction(pvListDataGetter, pv.PVListScrapperTypes)
}

func BuildPVDetailAction() ActionHandler {
	return detailResourceAction(pv.PVDetailScrapperTypes)
}
