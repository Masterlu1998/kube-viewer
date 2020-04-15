package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pvc"
	"github.com/Masterlu1998/kube-viewer/terminal/path"
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

func BuildPVCListAction(tree *path.TrieTree) {
	listResourceAction(pvcListDataGetter, tree, pvc.PVCListScrapperTypes)
}

func BuildPVCDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, pvc.PVCDetailScrapperTypes)
}
