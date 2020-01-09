package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
)

var (
	workloadTableHeader   = []string{"name", "namespace", "pods", "create time", "images"}
	workloadTableColWidth = []int{40, 22, 10, 30, 40}
)

func workloadListDataGetter(c common.KubernetesData) ([]string, [][]string, []int, error) {
	workloadSData, ok := c.([]workload.Info)
	if !ok {
		return nil, nil, nil, errors.New("convert to workload.Info failed")
	}

	newWorkloadTableData := make([][]string, 0)
	for _, wd := range workloadSData {
		newWorkloadTableData = append(newWorkloadTableData, []string{
			wd.Name,
			wd.Namespace,
			wd.PodsLive + "/" + wd.PodsTotal,
			wd.CreateTime,
			wd.Images,
		})
	}

	return workloadTableHeader, newWorkloadTableData, workloadTableColWidth, nil
}

func BuildDeploymentListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.DeploymentScrapperTypes)
}

func BuildStatefulSetListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.StatefulSetScrapperTypes)
}

func BuildDaemonSetListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.DaemonSetScrapperTypes)
}

func BuildReplicaSetListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.ReplicaSetScrapperTypes)
}

func BuildJobListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.JobScrapperTypes)
}

func BuildCronJobListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.CronJobScrapperTypes)
}
