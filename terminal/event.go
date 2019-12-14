package terminal

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	ui "github.com/gizak/termui/v3"
)

var (
	resourceTypesList = []string{
		workload.DeploymentResourceTypes,
		workload.StatefulSetResourceTypes,
		workload.DaemonSetResourceTypes,
		workload.ReplicaSetResourceTypes,
		workload.CronJobResourceTypes,
		workload.JobResourceTypes,
	}
)

type eventListener struct {
	ctx                context.Context
	tdb                *TerminalDashBoard
	cancelFunc         context.CancelFunc
	resourceTypesList  []string
	namespacesList     []string
	pathHandlerMap     map[string]handler
	resourceTypesIndex int
	namespacesIndex    int
	scrapperManagement *kScrapper.ScrapperManagement
}

type handler func(ctx context.Context, tdb *TerminalDashBoard, sm *kScrapper.ScrapperManagement, namespace string)

func newEventListener(ctx context.Context, tdb *TerminalDashBoard, cancel context.CancelFunc, sm *kScrapper.ScrapperManagement) *eventListener {
	return &eventListener{
		ctx:                ctx,
		tdb:                tdb,
		cancelFunc:         cancel,
		scrapperManagement: sm,
		resourceTypesList:  resourceTypesList,
		namespacesList:     []string{""},
		pathHandlerMap:     make(map[string]handler),
		resourceTypesIndex: 0,
		namespacesIndex:    0,
	}
}

func (el *eventListener) Register(path string, h handler) {
	el.pathHandlerMap[path] = h
}

func (el *eventListener) Listen() error {
	el.executeHandler("/" + el.resourceTypesList[el.resourceTypesIndex] + "/list")
	el.tdb.AddResourcePointer(0)
	el.syncNamespace()

	for {
		e := <-ui.PollEvents()
		switch e.ID {
		case "<Resize>":
			el.tdb.Resize()
		case "q":
			el.cancelFunc()
			return nil
		case "<Left>":
			el.tdb.MoveTabLeft()
			if el.namespacesIndex > 0 {
				el.namespacesIndex = el.namespacesIndex - 1
			}
			path := "/" + el.resourceTypesList[el.resourceTypesIndex] + "/list"
			el.executeHandler(path)
		case "<Right>":
			el.tdb.MoveTabRight()
			if el.namespacesIndex < len(el.namespacesList) {
				el.namespacesIndex = el.namespacesIndex + 1
			}
			path := "/" + el.resourceTypesList[el.resourceTypesIndex] + "/list"
			el.executeHandler(path)
		case "<Up>":
			el.tdb.RemoveResourcePointer(el.resourceTypesIndex)
			if el.resourceTypesIndex > 0 {
				el.resourceTypesIndex = el.resourceTypesIndex - 1
			}
			el.tdb.AddResourcePointer(el.resourceTypesIndex)
			path := "/" + el.resourceTypesList[el.resourceTypesIndex] + "/list"
			el.executeHandler(path)
		case "<Down>":
			el.tdb.RemoveResourcePointer(el.resourceTypesIndex)
			if el.resourceTypesIndex < len(el.resourceTypesList)-1 {
				el.resourceTypesIndex = el.resourceTypesIndex + 1
			}
			el.tdb.AddResourcePointer(el.resourceTypesIndex)
			path := "/" + el.resourceTypesList[el.resourceTypesIndex] + "/list"
			el.executeHandler(path)
		}

	}
}

func (el *eventListener) executeHandler(path string) {
	if handler, ok := el.pathHandlerMap[path]; ok {
		go handler(el.ctx, el.tdb, el.scrapperManagement, el.namespacesList[el.namespacesIndex])
	}

	return
}

func (el *eventListener) syncNamespace() {
	nc := el.scrapperManagement.ScrapperMap[namespace.NamespaceScrapperTypes]
	nc.StartScrapper(el.ctx, "")
	ns := <-nc.Watch()
	nc.StopResourceScrapper()
	namespaces := ns.([]string)
	el.namespacesList = namespaces
	el.tdb.NamespaceTab.TabNames = namespaces
	ui.Render(el.tdb.Grid)
}
