package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/configMap"
	"github.com/Masterlu1998/kube-viewer/terminal/path"
)

var (
	configMapTableHeader   = []string{"name", "namespace", "create time"}
	configMapTableColWidth = []int{35, 25, 20}
)

func configMapListDataGetter(c common.KubernetesData) ([]string, [][]string, []int, error) {
	configMapInfos, ok := c.([]configMap.Info)
	if !ok {
		return nil, nil, nil, errors.New("convert to configMap.Info failed")
	}

	newConfigMapTableData := make([][]string, 0)
	for _, cInfo := range configMapInfos {
		newConfigMapTableData = append(newConfigMapTableData, []string{
			cInfo.Name,
			cInfo.Namespace,
			cInfo.CreateTime,
		})
	}

	return configMapTableHeader, newConfigMapTableData, configMapTableColWidth, nil
}

func BuildConfigMapListAction(tree *path.TrieTree) {
	listResourceAction(configMapListDataGetter, tree, configMap.ConfigListMapScrapperTypes)
}

func BuildConfigMapDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, configMap.ConfigMapDetailScrapperTypes)
}
