package terminal

import (
	"fmt"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/configMap"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pv"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pvc"
	"github.com/Masterlu1998/kube-viewer/kScrapper/secret"
	"github.com/Masterlu1998/kube-viewer/kScrapper/service"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	ui "github.com/gizak/termui/v3"
)

var (
	resourceTableHeader  = [][]string{{"name", "namespace", "pods", "create time", "images"}}
	serviceTableHeader   = [][]string{{"name", "namespace", "clusterIP", "ports"}}
	configMapTableHeader = [][]string{{"name", "namespace", "create time"}}
	pvcTableHeader       = [][]string{{"name", "namespace", "status", "volume", "request", "limit", "accessMode", "storageClass", "create time"}}
	pvTableHeader        = [][]string{{"name", "capacity", "accessMode", "reclaim policy", "status", "storage class", "create time"}}
)

func (el *eventListener) pvGraphAction() {
	el.scrapperManagement.StopMainScrapper()
	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, pv.PVScrapperTypes, el.getCurrentNamespace())
	if err != nil {
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "PVAction"))
		return
	}

	for c := range el.scrapperManagement.GetSpecificScrapperCh(pv.PVScrapperTypes) {
		pvInfos, ok := c.([]pv.Info)
		if !ok {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, fmt.Sprintf("convert to pv.Info failed"), "PVAction"))
			return
		}

		newPVTableData := pvTableHeader
		for _, pvcInfo := range pvInfos {
			newPVTableData = append(newPVTableData, []string{
				pvcInfo.Name,
				pvcInfo.Capacity,
				pvcInfo.AccessMode,
				pvcInfo.ReclaimPolicy,
				pvcInfo.Status,
				pvcInfo.StorageClass,
				pvcInfo.CreateTime,
			})
		}
		el.tdb.ResourceTable.RefreshTableData(newPVTableData)
		ui.Render(el.tdb)
	}
}

func (el *eventListener) pvcGraphAction() {
	el.scrapperManagement.StopMainScrapper()
	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, pvc.PVCScrapperTypes, el.getCurrentNamespace())
	if err != nil {
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "PVCAction"))
		return
	}

	for c := range el.scrapperManagement.GetSpecificScrapperCh(pvc.PVCScrapperTypes) {
		pvcInfos, ok := c.([]pvc.Info)
		if !ok {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, fmt.Sprintf("convert to pvc.Info failed"), "PVCAction"))
			return
		}

		newPVCTableData := pvcTableHeader
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
		el.tdb.ResourceTable.RefreshTableData(newPVCTableData)
		ui.Render(el.tdb)
	}
}

func (el *eventListener) secretGraphAction() {
	el.scrapperManagement.StopMainScrapper()
	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, secret.SecretScrapperTypes, el.getCurrentNamespace())
	if err != nil {
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "SecretAction"))
		return
	}

	for c := range el.scrapperManagement.GetSpecificScrapperCh(secret.SecretScrapperTypes) {
		configMapInfos, ok := c.([]secret.Info)
		if !ok {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, fmt.Sprintf("convert to secret.Info failed"), "SecretAction"))
			return
		}

		newConfigMapTableData := configMapTableHeader
		for _, cInfo := range configMapInfos {
			newConfigMapTableData = append(newConfigMapTableData, []string{
				cInfo.Name,
				cInfo.Namespace,
				cInfo.CreateTime,
			})
		}
		el.tdb.ResourceTable.RefreshTableData(newConfigMapTableData)
		ui.Render(el.tdb)
	}
}

func (el *eventListener) configMapGraphAction() {
	el.scrapperManagement.StopMainScrapper()
	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, configMap.ConfigMapScrapperTypes, el.getCurrentNamespace())
	if err != nil {
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "ServiceAction"))
		return
	}

	for c := range el.scrapperManagement.GetSpecificScrapperCh(configMap.ConfigMapScrapperTypes) {
		configMapInfos, ok := c.([]configMap.Info)
		if !ok {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, fmt.Sprintf("convert to configMap.Info failed"), "ConfigMapAction"))
			return
		}

		newConfigMapTableData := configMapTableHeader
		for _, cInfo := range configMapInfos {
			newConfigMapTableData = append(newConfigMapTableData, []string{
				cInfo.Name,
				cInfo.Namespace,
				cInfo.CreateTime,
			})
		}
		el.tdb.ResourceTable.RefreshTableData(newConfigMapTableData)
		ui.Render(el.tdb)
	}
}

func (el *eventListener) serviceGraphAction() {
	el.scrapperManagement.StopMainScrapper()
	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, service.ServiceScrapperTypes, el.getCurrentNamespace())
	if err != nil {
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "ServiceAction"))
		return
	}

	for s := range el.scrapperManagement.GetSpecificScrapperCh(service.ServiceScrapperTypes) {
		newServiceTableData := serviceTableHeader
		serviceInfos, ok := s.([]service.Info)
		if !ok {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, fmt.Sprintf("convert to service.Info failed"), "ServiceAction"))
		}

		for _, serviceInfo := range serviceInfos {
			newServiceTableData = append(newServiceTableData, []string{
				serviceInfo.Name,
				serviceInfo.Namespace,
				serviceInfo.ClusterIP,
				serviceInfo.Port,
			})
		}
		el.tdb.ResourceTable.RefreshTableData(newServiceTableData)
		ui.Render(el.tdb)
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

		newWorkloadTableData := resourceTableHeader
		for _, wd := range workloadSData {
			newWorkloadTableData = append(newWorkloadTableData, []string{
				wd.Name,
				wd.Namespace,
				wd.PodsLive + "/" + wd.PodsTotal,
				wd.CreateTime,
				wd.Images,
			})
		}
		el.tdb.ResourceTable.RefreshTableData(newWorkloadTableData)
		ui.Render(el.tdb)
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
			el.tdb.NamespaceTab.RefreshNamespace(namespaces)
		}
		ui.Render(el.tdb)
	}
}

func (el *eventListener) collectDebugMessage() {
	for m := range el.debugCollector.GetDebugMessageCh() {
		el.tdb.Console.ShowDebugInfo(m)
		ui.Render(el.tdb)
	}
}

func (el *eventListener) leftKeyboardAction() {
	el.tdb.NamespaceTab.FocusLeft()
	if el.namespacesIndex > 0 {
		el.namespacesIndex = el.namespacesIndex - 1
	}
	el.scrapperManagement.ResetNamespace(el.getCurrentNamespace())
	ui.Render(el.tdb)
}

func (el *eventListener) rightKeyboardAction() {
	el.tdb.NamespaceTab.FocusRight()
	if el.namespacesIndex < len(el.namespacesList)-1 {
		el.namespacesIndex = el.namespacesIndex + 1
	}
	el.scrapperManagement.ResetNamespace(el.getCurrentNamespace())
	ui.Render(el.tdb)
}

func (el *eventListener) upKeyboardAction() {
	el.tdb.Menu.ScrollUp()
	ui.Render(el.tdb)
}

func (el *eventListener) downKeyboardAction() {
	el.tdb.Menu.ScrollDown()
	ui.Render(el.tdb)
}

func (el *eventListener) enterKeyboardAction() {
	path := el.tdb.Menu.Enter()
	el.executeHandler(path)
}
