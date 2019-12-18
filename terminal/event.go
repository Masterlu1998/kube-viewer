package terminal

import (
	"context"
	"fmt"

	// "fmt"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/configMap"
	"github.com/Masterlu1998/kube-viewer/kScrapper/service"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
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
	keyboardActionTypes     = "keyboard"
	namespaceActionTypes    = "namespace"
	debugMessageActionTypes = "debug"
)

type eventListener struct {
	ctx                context.Context
	tdb                *component.TerminalDashBoard
	cancelFunc         context.CancelFunc
	resourceTypesList  []string
	namespacesList     []string
	pathHandlerMap     map[string]handler
	namespacesIndex    int
	scrapperManagement *kScrapper.ScrapperManagement
	debugCollector     *debug.DebugCollector
}

type handler func()

func newEventListener(ctx context.Context, tdb *component.TerminalDashBoard, cancel context.CancelFunc,
	sm *kScrapper.ScrapperManagement,
	dc *debug.DebugCollector) *eventListener {
	return &eventListener{
		ctx:                ctx,
		tdb:                tdb,
		cancelFunc:         cancel,
		scrapperManagement: sm,
		resourceTypesList:  resourceTypesList,
		namespacesList:     []string{""},
		pathHandlerMap:     make(map[string]handler),
		namespacesIndex:    0,
		debugCollector:     dc,
	}
}

func (el *eventListener) Register() {
	el.pathHandlerMap = map[string]handler{
		"/" + keyboardActionTypes + "/left":  el.leftKeyboardAction,
		"/" + keyboardActionTypes + "/right": el.rightKeyboardAction,
		"/" + keyboardActionTypes + "/up":    el.upKeyboardAction,
		"/" + keyboardActionTypes + "/down":  el.downKeyboardAction,
		"/" + keyboardActionTypes + "/enter": el.enterKeyboardAction,

		"/" + namespaceActionTypes + "/sync":              el.syncNamespaceAction,
		"/" + workload.DeploymentResourceTypes + "/list":  el.deploymentGraphAction,
		"/" + workload.StatefulSetResourceTypes + "/list": el.statefulSetGraphAction,
		"/" + workload.DaemonSetResourceTypes + "/list":   el.daemonSetGraphAction,
		"/" + workload.ReplicaSetResourceTypes + "/list":  el.replicaSetGraphAction,
		"/" + workload.CronJobResourceTypes + "/list":     el.cronJobGraphAction,
		"/" + workload.JobResourceTypes + "/list":         el.jobGraphAction,
		"/" + service.ServiceResourceTypes + "/list":      el.serviceGraphAction,
		"/" + configMap.ConfigMapResourceTypes + "/list":  el.configMapAction,
		"/" + debugMessageActionTypes + "/collect":        el.collectDebugMessage,
	}
}

func (el *eventListener) Listen() error {
	el.executeHandler("/" + debugMessageActionTypes + "/collect")
	el.executeHandler("/" + namespaceActionTypes + "/sync")

	for {
		e := <-ui.PollEvents()
		switch e.ID {
		case "<Resize>":
			el.tdb.Resize()
		case "q":
			el.cancelFunc()
			return nil
		case "<Enter>":
			path := "/" + keyboardActionTypes + "/enter"
			el.executeHandler(path)
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
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Info, fmt.Sprintf("excute path: %s", path), "eventListener"))
		return
	}
	el.debugCollector.Collect(debug.NewDebugMessage(debug.Warn, fmt.Sprintf("no action match path: %s", path), "eventListener"))
}

func (el *eventListener) getCurrentNamespace() string {
	return el.namespacesList[el.namespacesIndex]
}
