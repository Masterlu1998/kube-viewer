package terminal

import (
	"fmt"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/service"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	ui "github.com/gizak/termui/v3"
)

var (
	resourceTableHeader = [][]string{{"name", "namespace", "pods", "create time", "images"}}
	serviceTableHeader  = [][]string{{"name", "namespace", "clusterIP", "ports"}}
)

func (el *eventListener) serviceGraphAction() {
	el.scrapperManagement.StopMainScrapper()
	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, service.ServiceScrapperTypes, el.getCurrentNamespace())
	if err != nil {
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "ServiceAction"))
		return
	}

	defer el.debugCollector.Collect(debug.NewDebugMessage(debug.Info, "ServiceAction stop", "ServiceAction"))
	for s := range el.scrapperManagement.GetSpecificScrapperCh(service.ServiceScrapperTypes) {
		el.tdb.ResourceTable.Rows = serviceTableHeader
		sInfos, ok := s.([]service.Info)
		if !ok {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, fmt.Sprintf("convert to service.Info failed"), "ServiceAction"))
		}

		for _, s := range sInfos {
			var serviceContent []string
			serviceContent = append(serviceContent,
				s.Name,
				s.Namespace,
				s.ClusterIP,
				s.Port,
			)
			el.tdb.ResourceTable.Rows = append(el.tdb.ResourceTable.Rows, serviceContent)
		}
		ui.Render(el.tdb.Grid)
	}
}

func (el *eventListener) deploymentGraphAction() {
	el.workloadGraphAction(workload.DeploymentScrapperTypes)
}

func (el *eventListener) statefulSetGraphAction() {
	el.workloadGraphAction(workload.StatefulSetScrapperTypes)
}

func (el *eventListener) daemonSetGraphAction() {
	el.workloadGraphAction(workload.DaemonSetScrapperTypes)

}

func (el *eventListener) replicaSetGraphAction() {
	el.workloadGraphAction(workload.ReplicaSetScrapperTypes)
}

func (el *eventListener) cronJobGraphAction() {
	el.workloadGraphAction(workload.CronJobScrapperTypes)
}

func (el *eventListener) jobGraphAction() {
	el.workloadGraphAction(workload.JobScrapperTypes)
}

func (el *eventListener) workloadGraphAction(scrapperType string) {
	el.scrapperManagement.StopMainScrapper()
	el.debugCollector.Collect(debug.NewDebugMessage(debug.Info, "start scrapper", scrapperType))

	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, scrapperType, el.getCurrentNamespace())
	if err != nil {
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "resourceType"))
		return
	}

	t := el.tdb.ResourceTable
	for d := range el.scrapperManagement.GetSpecificScrapperCh(scrapperType) {
		workloadSData, ok := d.([]workload.WorkloadInfo)
		if !ok {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Error,
				"convert to type []workload.WorkloadInfo failed", "workloadAction"))
			continue
		}

		if len(workloadSData) == 0 {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Warn,
				"workload list is empty", "workloadAction"))
		}

		t.Rows = resourceTableHeader
		for _, wd := range workloadSData {
			var deploymentContent []string
			deploymentContent = append(deploymentContent,
				wd.Name,
				wd.Namespace,
				wd.PodsLive+"/"+wd.PodsTotal,
				wd.CreateTime,
				wd.Images)
			t.Rows = append(t.Rows, deploymentContent)
		}
		ui.Render(el.tdb.Grid)
	}
}

func (el *eventListener) syncNamespaceAction() {
	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, namespace.NamespaceScrapperTypes, "")
	if err != nil {
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "namespaceAction"))
		return
	}

	for {
		select {
		case <-el.ctx.Done():
			return
		case ns := <-el.scrapperManagement.GetSpecificScrapperCh(namespace.NamespaceScrapperTypes):
			namespaces := []string{""}
			namespaces = append(namespaces, ns.([]string)...)
			el.namespacesList = namespaces
			el.tdb.NamespaceTab.TabNames = namespaces
		}
		ui.Render(el.tdb.Grid)
	}
}

func (el *eventListener) collectDebugMessage() {
	for m := range el.debugCollector.GetDebugMessageCh() {
		el.tdb.ShowDebugInfo(m)
		ui.Render(el.tdb.Grid)
	}
}

func (el *eventListener) leftKeyboardAction() {
	el.tdb.MoveTabLeft()
	if el.namespacesIndex > 0 {
		el.namespacesIndex = el.namespacesIndex - 1
	}
	el.scrapperManagement.ResetNamespace(el.getCurrentNamespace())
	ui.Render(el.tdb.Grid)
}

func (el *eventListener) rightKeyboardAction() {
	el.tdb.MoveTabRight()
	if el.namespacesIndex < len(el.namespacesList)-1 {
		el.namespacesIndex = el.namespacesIndex + 1
	}
	el.scrapperManagement.ResetNamespace(el.getCurrentNamespace())
	ui.Render(el.tdb.Grid)
}

func (el *eventListener) upKeyboardAction() {
	el.tdb.SelectUp()
	ui.Render(el.tdb.Grid)
}

func (el *eventListener) downKeyboardAction() {
	el.tdb.SelectDown()
	ui.Render(el.tdb.Grid)
}

func (el *eventListener) enterKeyboardAction() {
	path := el.tdb.Enter()
	el.executeHandler(path)
}
