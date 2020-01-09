package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/configMap"
)

var (
	configMapTableHeader   = []string{"name", "namespace", "create time"}
	configMapTableColWidth = []int{40, 25, 40}
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

func BuildConfigMapListAction() ActionHandler {
	return listResourceAction(configMapListDataGetter, configMap.ConfigMapScrapperTypes)
}
