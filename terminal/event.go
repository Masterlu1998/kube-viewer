package terminal

import (
	"context"
	"fmt"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/configMap"
	"github.com/Masterlu1998/kube-viewer/kScrapper/node"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pv"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pvc"
	"github.com/Masterlu1998/kube-viewer/kScrapper/secret"
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
	terminalDashBoard  *component.TerminalDashBoard
	cancelFunc         context.CancelFunc
	resourceTypesList  []string
	namespacesList     []string
	pathHandlerMap     map[string]handler
	namespacesIndex    int
	scrapperManagement *kScrapper.ScrapperManagement
	debugCollector     *debug.DebugCollector
	currentPanel       component.PanelTypes
}

type handler func()

func newEventListener(ctx context.Context, tdb *component.TerminalDashBoard, cancel context.CancelFunc,
	sm *kScrapper.ScrapperManagement,
	dc *debug.DebugCollector) *eventListener {
	return &eventListener{
		ctx:                ctx,
		terminalDashBoard:  tdb,
		cancelFunc:         cancel,
		scrapperManagement: sm,
		resourceTypesList:  resourceTypesList,
		namespacesList:     []string{""},
		pathHandlerMap:     make(map[string]handler),
		namespacesIndex:    0,
		debugCollector:     dc,
		currentPanel:       component.MenuPanel,
	}
}

func (el *eventListener) Register() {
	el.pathHandlerMap = map[string]handler{
		"/" + keyboardActionTypes + "/left":                                             el.leftKeyboardAction,
		"/" + keyboardActionTypes + "/right":                                            el.rightKeyboardAction,
		"/" + string(component.MenuPanel) + "/" + keyboardActionTypes + "/up":           el.upMenuKeyboardAction,
		"/" + string(component.MenuPanel) + "/" + keyboardActionTypes + "/down":         el.downMenuKeyboardAction,
		"/" + string(component.MenuPanel) + "/" + keyboardActionTypes + "/enter":        el.enterMenuKeyboardAction,
		"/" + string(component.ResourceListPanel) + "/" + keyboardActionTypes + "/up":   el.upResourceListKeyboardAction,
		"/" + string(component.ResourceListPanel) + "/" + keyboardActionTypes + "/down": el.downResourceListKeyboardAction,
		// "/" + string(component.ResourceListPanel) + "/" + keyboardActionTypes + "/enter": el.enterMenuKeyboardAction,
		"/" + keyboardActionTypes + "/tab": el.tabKeyboardAction,

		"/" + namespaceActionTypes + "/sync":              el.syncNamespaceAction,
		"/" + workload.DeploymentResourceTypes + "/list":  el.deploymentGraphAction,
		"/" + workload.StatefulSetResourceTypes + "/list": el.statefulSetGraphAction,
		"/" + workload.DaemonSetResourceTypes + "/list":   el.daemonSetGraphAction,
		"/" + workload.ReplicaSetResourceTypes + "/list":  el.replicaSetGraphAction,
		"/" + workload.CronJobResourceTypes + "/list":     el.cronJobGraphAction,
		"/" + workload.JobResourceTypes + "/list":         el.jobGraphAction,
		"/" + service.ServiceResourceTypes + "/list":      el.serviceGraphAction,
		"/" + configMap.ConfigMapResourceTypes + "/list":  el.configMapGraphAction,
		"/" + secret.SecretResourceTypes + "/list":        el.secretGraphAction,
		"/" + pvc.PVCResourceTypes + "/list":              el.pvcGraphAction,
		"/" + pv.PVResourceTypes + "/list":                el.pvGraphAction,
		"/" + node.NodeResourceTypes + "/list":            el.nodeGraphAction,

		"/" + debugMessageActionTypes + "/collect": el.collectDebugMessage,
	}
}

func (el *eventListener) Listen() error {
	el.executeHandler("/" + debugMessageActionTypes + "/collect")
	el.executeHandler("/" + namespaceActionTypes + "/sync")

	for {
		e := <-ui.PollEvents()
		switch e.ID {
		case "<Resize>":
			el.terminalDashBoard.Resize()
		case "q":
			el.cancelFunc()
			return nil
		case "<Tab>":
			path := "/" + keyboardActionTypes + "/tab"
			el.executeHandler(path)
		case "<Enter>":
			path := "/" + string(el.currentPanel) + "/" + keyboardActionTypes + "/enter"
			el.executeHandler(path)
		case "<Left>":
			path := "/" + keyboardActionTypes + "/left"
			el.executeHandler(path)
		case "<Right>":
			path := "/" + keyboardActionTypes + "/right"
			el.executeHandler(path)
		case "<Up>":
			path := "/" + string(el.currentPanel) + "/" + keyboardActionTypes + "/up"
			el.executeHandler(path)
		case "<Down>":
			path := "/" + string(el.currentPanel) + "/" + keyboardActionTypes + "/down"
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
