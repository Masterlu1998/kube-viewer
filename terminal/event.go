package terminal

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
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

const (
	keyboardActionTypes  = "keyboard"
	namespaceActionTypes = "namespace"
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

type handler func()

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

func (el *eventListener) Register() {
	el.pathHandlerMap = map[string]handler{
		"/" + keyboardActionTypes + "/left":              el.leftKeyboardAction,
		"/" + keyboardActionTypes + "/right":             el.rightKeyboardAction,
		"/" + keyboardActionTypes + "/up":                el.upKeyboardAction,
		"/" + keyboardActionTypes + "/down":              el.downKeyboardAction,
		"/" + namespaceActionTypes + "/sync":             el.syncNamespaceAction,
		"/" + workload.DeploymentResourceTypes + "/list": el.deploymentGraphAction,
	}
}

func (el *eventListener) Listen() error {
	el.executeHandler("/" + namespaceActionTypes + "/sync")
	el.executeHandler("/" + el.getCurrentResourceType() + "/list")
	el.tdb.AddResourcePointer(0)

	for {
		e := <-ui.PollEvents()
		switch e.ID {
		case "<Resize>":
			el.tdb.Resize()
		case "q":
			el.cancelFunc()
			return nil
		case "<Left>":
			path := "/" + keyboardActionTypes + "/left"
			el.executeHandler(path)
		case "<Right>":
			path := "/" + keyboardActionTypes + "/right"
			el.executeHandler(path)
		case "<Up>":
			path := "/" + keyboardActionTypes + "/up"
			el.executeHandler(path)
		case "<Down>":
			path := "/" + keyboardActionTypes + "/down"
			el.executeHandler(path)
		}
	}
}

func (el *eventListener) executeHandler(path string) {
	if handler, ok := el.pathHandlerMap[path]; ok {
		go handler()
	}

	return
}

func (el *eventListener) getCurrentNamespace() string {
	return el.namespacesList[el.namespacesIndex]
}

func (el *eventListener) getCurrentResourceType() string {
	return el.resourceTypesList[el.namespacesIndex]
}
