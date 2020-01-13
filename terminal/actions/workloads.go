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
	return listResourceAction(workloadListDataGetter, workload.DeploymentListScrapperTypes)
}

func BuildDeploymentDetailAction() ActionHandler {
	return detailResourceAction(workload.DeploymentDetailScrapperTypes)
}

func BuildStatefulSetListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.StatefulSetListScrapperTypes)
}

func BuildStatefulSetDetailAction() ActionHandler {
	return detailResourceAction(workload.StatefulSetDetailScrapperTypes)
}

func BuildDaemonSetListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.DaemonSetListScrapperTypes)
}

func BuildDaemonSetDetailAction() ActionHandler {
	return detailResourceAction(workload.DaemonSetDetailScrapperTypes)
}

func BuildReplicaSetListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.ReplicaSetListScrapperTypes)
}

func BuildReplicaSetDetailAction() ActionHandler {
	return detailResourceAction(workload.ReplicaSetDetailScrapperTypes)
}

func BuildJobListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.JobListScrapperTypes)
}

func BuildJobDetailAction() ActionHandler {
	return detailResourceAction(workload.JobDetailScrapperTypes)
}

func BuildCronJobListAction() ActionHandler {
	return listResourceAction(workloadListDataGetter, workload.CronJobListScrapperTypes)
}

func BuildCronJobDetailAction() ActionHandler {
	return detailResourceAction(workload.CronJobDetailScrapperTypes)
}
