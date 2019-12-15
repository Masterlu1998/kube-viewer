package terminal

import (
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	ui "github.com/gizak/termui/v3"
)

var resourceTableHeader = [][]string{{"name", "namespace", "pods", "create time", "images"}}

func (el *eventListener) workloadGraphAction() {
	err := el.scrapperManagement.StartSpecificScrapper(el.ctx, el.getCurrentScrapperType(), el.getCurrentNamespace())
	if err != nil {
		return
	}

	t := el.tdb.ResourceTable
	for d := range el.scrapperManagement.GetSpecificScrapperCh(el.getCurrentScrapperType()) {
		workloadSData, ok := d.([]workload.WorkloadInfo)
		if !ok {
			continue
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
		return
	}

	go func() {
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
	}()
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
	el.tdb.RemoveResourcePointer(el.resourceTypesIndex)
	if el.resourceTypesIndex > 0 {
		el.scrapperManagement.StopSpecificScrapper(el.getCurrentScrapperType())
		el.resourceTypesIndex = el.resourceTypesIndex - 1
	}
	el.tdb.AddResourcePointer(el.resourceTypesIndex)
	path := "/" + workloadActionTypes + "/list"
	el.executeHandler(path)
	ui.Render(el.tdb.Grid)
}

func (el *eventListener) downKeyboardAction() {
	el.tdb.RemoveResourcePointer(el.resourceTypesIndex)
	if el.resourceTypesIndex < len(el.resourceTypesList)-1 {
		el.scrapperManagement.StopSpecificScrapper(el.getCurrentScrapperType())
		el.resourceTypesIndex = el.resourceTypesIndex + 1
	}
	el.tdb.AddResourcePointer(el.resourceTypesIndex)
	path := "/" + workloadActionTypes + "/list"
	el.executeHandler(path)
	ui.Render(el.tdb.Grid)
}
