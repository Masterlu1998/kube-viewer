package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	"github.com/Masterlu1998/kube-viewer/terminal/path"
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

func BuildDeploymentListAction(tree *path.TrieTree) {
	listResourceAction(workloadListDataGetter, tree, workload.DeploymentListScrapperTypes)
}

func BuildDeploymentDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, workload.DeploymentDetailScrapperTypes)
}

func BuildStatefulSetListAction(tree *path.TrieTree) {
	listResourceAction(workloadListDataGetter, tree, workload.StatefulSetListScrapperTypes)
}

func BuildStatefulSetDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, workload.StatefulSetDetailScrapperTypes)
}

func BuildDaemonSetListAction(tree *path.TrieTree) {
	listResourceAction(workloadListDataGetter, tree, workload.DaemonSetListScrapperTypes)
}

func BuildDaemonSetDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, workload.DaemonSetDetailScrapperTypes)
}

func BuildReplicaSetListAction(tree *path.TrieTree) {
	listResourceAction(workloadListDataGetter, tree, workload.ReplicaSetListScrapperTypes)
}

func BuildReplicaSetDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, workload.ReplicaSetDetailScrapperTypes)
}

func BuildJobListAction(tree *path.TrieTree) {
	listResourceAction(workloadListDataGetter, tree, workload.JobListScrapperTypes)
}

func BuildJobDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, workload.JobDetailScrapperTypes)
}

func BuildCronJobListAction(tree *path.TrieTree) {
	listResourceAction(workloadListDataGetter, tree, workload.CronJobListScrapperTypes)
}

func BuildCronJobDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, workload.CronJobDetailScrapperTypes)
}
