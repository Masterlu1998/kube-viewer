package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/service"
)

var (
	serviceTableHeader   = []string{"name", "namespace", "clusterIP", "ports"}
	serviceTableColWidth = []int{30, 25, 20, 50}
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

func BuildServiceListAction() ActionHandler {
	return listResourceAction(serviceListDataGetter, service.ServiceListScrapperTypes)
}

func BuildServiceDetailAction() ActionHandler {
	return detailResourceAction(service.ServiceDetailScrapperTypes)
}
