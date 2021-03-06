package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/service"
	"github.com/Masterlu1998/kube-viewer/terminal/path"
)

var (
	serviceTableHeader   = []string{"name", "namespace", "clusterIP", "ports"}
	serviceTableColWidth = []int{20, 17, 15, 28}
)

func serviceListDataGetter(c common.KubernetesData) ([]string, [][]string, []int, error) {
	newServiceTableData := make([][]string, 0)
	serviceInfos, ok := c.([]service.Info)
	if !ok {
		return nil, nil, nil, errors.New("convert to service.Info failed")
	}

	for _, serviceInfo := range serviceInfos {
		newServiceTableData = append(newServiceTableData, []string{
			serviceInfo.Name,
			serviceInfo.Namespace,
			serviceInfo.ClusterIP,
			serviceInfo.Port,
		})
	}

	return serviceTableHeader, newServiceTableData, serviceTableColWidth, nil
}

func BuildServiceListAction(tree *path.TrieTree) {
	listResourceAction(serviceListDataGetter, tree, service.ServiceListScrapperTypes)
}

func BuildServiceDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, service.ServiceDetailScrapperTypes)
}
