package terminal

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/resource"
	ui "github.com/gizak/termui/v3"
)

var (
	resourceTypesList = []resource.ResourceTypes{
		resource.DeploymentResourceTypes,
		resource.StatefulSetResourceTypes,
		resource.DaemonSetResourceTypes,
		resource.ReplicaSetResourceTypes,
		resource.CronJobResourceTypes,
		resource.JobResourceTypes,
	}
)

type eventListener struct {
	ctx                context.Context
	tdb                *TerminalDashBoard
	cancelFunc         context.CancelFunc
	resourceTypesList  []resource.ResourceTypes
	namespacesList     []string
	pathHandlerMap     map[string]handler
	resourceTypesIndex int
	namespacesIndex    int
	resourceScrapper   *resource.ResourceScrapper
}

type handler func(ctx context.Context, tdb *TerminalDashBoard, s *resource.ResourceScrapper, workloadTypes resource.ResourceTypes)

func newEventListener(ctx context.Context, tdb *TerminalDashBoard, cancel context.CancelFunc, s *kScrapper.ScrapperManagement) *eventListener {
	return &eventListener{
		ctx:                ctx,
		tdb:                tdb,
		cancelFunc:         cancel,
		resourceScrapper:   s.ResourceScrapper,
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
	go el.pathHandlerMap["/workload/list"](el.ctx, el.tdb, el.resourceScrapper, el.resourceTypesList[el.resourceTypesIndex])
	el.tdb.AddResourcePointer(0)
	for {
		e := <-ui.PollEvents()
		switch e.ID {
		case "<Resize>":
			el.tdb.Resize()
		case "q":
			el.cancelFunc()
			return nil
		case "<Up>":
			el.tdb.RemoveResourcePointer(el.resourceTypesIndex)
			if el.resourceTypesIndex > 0 {
				el.resourceTypesIndex = el.resourceTypesIndex - 1
			}
			el.tdb.AddResourcePointer(el.resourceTypesIndex)
			path := "/workload/list"
			go el.pathHandlerMap[path](el.ctx, el.tdb, el.resourceScrapper, el.resourceTypesList[el.resourceTypesIndex])
		case "<Down>":
			el.tdb.RemoveResourcePointer(el.resourceTypesIndex)
			if el.resourceTypesIndex < len(el.resourceTypesList)-1 {
				el.resourceTypesIndex = el.resourceTypesIndex + 1
			}
			el.tdb.AddResourcePointer(el.resourceTypesIndex)
			path := "/workload/list"
			go el.pathHandlerMap[path](el.ctx, el.tdb, el.resourceScrapper, el.resourceTypesList[el.resourceTypesIndex])
		}

	}
}
