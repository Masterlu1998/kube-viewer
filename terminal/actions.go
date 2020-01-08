package terminal

import (
	"fmt"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/configMap"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/node"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pv"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pvc"
	"github.com/Masterlu1998/kube-viewer/kScrapper/secret"
	"github.com/Masterlu1998/kube-viewer/kScrapper/service"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	ui "github.com/gizak/termui/v3"
)

var (
	workloadTableHeader    = []string{"name", "namespace", "pods", "create time", "images"}
	workloadTableColWidth  = []int{40, 22, 10, 30, 40}
	serviceTableHeader     = []string{"name", "namespace", "clusterIP", "ports"}
	serviceTableColWidth   = []int{30, 25, 20, 50}
	configMapTableHeader   = []string{"name", "namespace", "create time"}
	configMapTableColWidth = []int{40, 25, 40}
	secretTableHeader      = []string{"name", "namespace", "create time"}
	secretTableColWidth    = []int{40, 25, 40}
	pvcTableHeader         = []string{"name", "namespace", "status", "volume", "request", "limit", "accessMode", "storageClass", "create time"}
	pvcTableColWidth       = []int{20, 20, 20, 20, 20, 20, 20, 20, 20}
	pvTableHeader          = []string{"name", "capacity", "accessMode", "reclaim policy", "status", "storage class", "create time"}
	pvTableColWidth        = []int{15, 10, 20, 20, 15, 20, 30}
	nodeTableHeader        = []string{"name", "status", "roles", "address", "OSImage"}
	nodeTableColWidth      = []int{20, 20, 15, 30, 20}
)

func (el *eventListener) nodeGraphAction() {
	el.scrapperManagement.StopMainScrapper()
	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, node.NodeScrapperTypes, el.getCurrentNamespace())
	if err != nil {
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "NodeAction"))
		return
	}

	for c := range el.scrapperManagement.GetSpecificScrapperCh(node.NodeScrapperTypes) {
		nodeInfos, ok := c.([]node.Info)
		if !ok {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, fmt.Sprintf("convert to pv.Info failed"), "PVAction"))
			return
		}

		newNodeTableData := make([][]string, 0)
		for _, nodeInfo := range nodeInfos {
			newNodeTableData = append(newNodeTableData, []string{
				nodeInfo.Name,
				nodeInfo.Status,
				nodeInfo.Roles,
				nodeInfo.Address,
				nodeInfo.OSImage,
			})
		}
		el.terminalDashBoard.ResourcePanel.RefreshPanelData(nodeTableHeader, newNodeTableData, nodeTableColWidth)
		ui.Render(el.terminalDashBoard)
	}
}

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

		newPVTableData := make([][]string, 0)
		for _, pvInfo := range pvInfos {
			newPVTableData = append(newPVTableData, []string{
				pvInfo.Name,
				pvInfo.Capacity,
				pvInfo.AccessMode,
				pvInfo.ReclaimPolicy,
				pvInfo.Status,
				pvInfo.StorageClass,
				pvInfo.CreateTime,
			})
		}
		el.terminalDashBoard.ResourcePanel.RefreshPanelData(pvTableHeader, newPVTableData, pvTableColWidth)
		ui.Render(el.terminalDashBoard)
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
		el.terminalDashBoard.ResourcePanel.RefreshPanelData(pvcTableHeader, newPVCTableData, pvcTableColWidth)
		ui.Render(el.terminalDashBoard)
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
		secretInfos, ok := c.([]secret.Info)
		if !ok {
			el.debugCollector.Collect(debug.NewDebugMessage(debug.Error, fmt.Sprintf("convert to secret.Info failed"), "SecretAction"))
			return
		}

		newSecretTableData := make([][]string, 0)
		for _, cInfo := range secretInfos {
			newSecretTableData = append(newSecretTableData, []string{
				cInfo.Name,
				cInfo.Namespace,
				cInfo.CreateTime,
			})
		}
		el.terminalDashBoard.ResourcePanel.RefreshPanelData(secretTableHeader, newSecretTableData, secretTableColWidth)
		ui.Render(el.terminalDashBoard)
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

		newConfigMapTableData := make([][]string, 0)
		for _, cInfo := range configMapInfos {
			newConfigMapTableData = append(newConfigMapTableData, []string{
				cInfo.Name,
				cInfo.Namespace,
				cInfo.CreateTime,
			})
		}
		el.terminalDashBoard.ResourcePanel.RefreshPanelData(configMapTableHeader, newConfigMapTableData, configMapTableColWidth)
		ui.Render(el.terminalDashBoard)
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
		newServiceTableData := make([][]string, 0)
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
		el.terminalDashBoard.ResourcePanel.RefreshPanelData(serviceTableHeader, newServiceTableData, serviceTableColWidth)
		ui.Render(el.terminalDashBoard)
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
		el.terminalDashBoard.ResourcePanel.RefreshPanelData(workloadTableHeader, newWorkloadTableData, workloadTableColWidth)
		ui.Render(el.terminalDashBoard)
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
			el.terminalDashBoard.NamespaceTab.RefreshNamespace(namespaces)
		}
		ui.Render(el.terminalDashBoard)
	}
}

func (el *eventListener) collectDebugMessage() {
	for m := range el.debugCollector.GetDebugMessageCh() {
		el.terminalDashBoard.Console.ShowDebugInfo(m)
		ui.Render(el.terminalDashBoard)
	}
}

func (el *eventListener) tabKeyboardAction() {
	el.terminalDashBoard.SwitchNextPanel()
	el.currentPanel = el.terminalDashBoard.GetCurrentPanelTypes()
}

func (el *eventListener) leftKeyboardAction() {
	el.terminalDashBoard.NamespaceTab.FocusLeft()
	if el.namespacesIndex > 0 {
		el.namespacesIndex = el.namespacesIndex - 1
	}
	el.scrapperManagement.ResetNamespace(el.getCurrentNamespace())
	ui.Render(el.terminalDashBoard)
}

func (el *eventListener) rightKeyboardAction() {
	el.terminalDashBoard.NamespaceTab.FocusRight()
	if el.namespacesIndex < len(el.namespacesList)-1 {
		el.namespacesIndex = el.namespacesIndex + 1
	}
	el.scrapperManagement.ResetNamespace(el.getCurrentNamespace())
	ui.Render(el.terminalDashBoard)
}

func (el *eventListener) upMenuKeyboardAction() {
	el.terminalDashBoard.Menu.ScrollUp()
	ui.Render(el.terminalDashBoard)
}

func (el *eventListener) downMenuKeyboardAction() {
	el.terminalDashBoard.Menu.ScrollDown()
	ui.Render(el.terminalDashBoard)
}

func (el *eventListener) enterMenuKeyboardAction() {
	path := el.terminalDashBoard.Menu.Enter()
	el.executeHandler(path)
}

func (el *eventListener) upResourceListKeyboardAction() {
	el.terminalDashBoard.ResourcePanel.ScrollUp()
	ui.Render(el.terminalDashBoard)
}

func (el *eventListener) downResourceListKeyboardAction() {
	el.terminalDashBoard.ResourcePanel.ScrollDown()
	ui.Render(el.terminalDashBoard)
}

// func (el *eventListener) enterResourceListKeyboardAction() {
// 	path := el.terminalDashBoard.Menu.Enter()
// 	el.executeHandler(path)
// }
