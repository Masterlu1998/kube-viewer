package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pvc"
)

var (
	pvcTableHeader   = []string{"name", "namespace", "status", "volume", "request", "limit", "accessMode", "storageClass", "create time"}
	pvcTableColWidth = []int{20, 20, 20, 20, 20, 20, 20, 20, 20}
)

func pvcListDataGetter(c common.KubernetesData) ([]string, [][]string, []int, error) {
	pvcInfos, ok := c.([]pvc.Info)
	if !ok {
		return nil, nil, nil, errors.New("convert to node.Info failed")
	}

	newPVCTableData := make([][]string, 0)
	for _, pvcInfo := range pvcInfos {
		newPVCTableData = append(newPVCTableData, []string{
			pvcInfo.Name,
			pvcInfo.Namespace,
			pvcInfo.Status,
			pvcInfo.Volume,
			pvcInfo.Request,
			pvcInfo.Limit,
			pvcInfo.AccessMode,
			pvcInfo.StorageClass,
			pvcInfo.CreateTime,
		})
	}

	return pvcTableHeader, newPVCTableData, pvcTableColWidth, nil
}

func BuildPVCListAction() ActionHandler {
	return listResourceAction(pvcListDataGetter, pvc.PVCListScrapperTypes)
}

func BuildPVCDetailAction() ActionHandler {
	return detailResourceAction(pvc.PVCDetailScrapperTypes)
}
