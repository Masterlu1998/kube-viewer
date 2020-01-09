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
	"github.com/Masterlu1998/kube-viewer/terminal/actions"
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
	pathHandlerMap     map[string]actions.ActionHandler
	scrapperManagement *kScrapper.ScrapperManagement
	debugCollector     *debug.DebugCollector
}

func newEventListener(ctx context.Context, tdb *component.TerminalDashBoard, cancel context.CancelFunc,
	sm *kScrapper.ScrapperManagement,
	dc *debug.DebugCollector) *eventListener {
	return &eventListener{
		ctx:                ctx,
		terminalDashBoard:  tdb,
		cancelFunc:         cancel,
		scrapperManagement: sm,
		resourceTypesList:  resourceTypesList,
		pathHandlerMap:     make(map[string]actions.ActionHandler),
		debugCollector:     dc,
	}
}

func (el *eventListener) Register() {
	el.pathHandlerMap = map[string]actions.ActionHandler{
		"/" + keyboardActionTypes + "/left":                                             actions.BuildLeftKeyboardAction(),
		"/" + keyboardActionTypes + "/right":                                            actions.BuildRightKeyboardAction(),
		"/" + string(component.MenuPanel) + "/" + keyboardActionTypes + "/up":           actions.BuildUpMenuKeyboardAction(),
		"/" + string(component.MenuPanel) + "/" + keyboardActionTypes + "/down":         actions.BuildDownMenuKeyboardAction(),
		"/" + string(component.ResourceListPanel) + "/" + keyboardActionTypes + "/up":   actions.BuildUpResourceListKeyboardAction(),
		"/" + string(component.ResourceListPanel) + "/" + keyboardActionTypes + "/down": actions.BuildDownResourceListKeyboardAction(),
		// "/" + string(component.ResourceListPanel) + "/" + keyboardActionTypes + "/enter": el.enterResourceListKeyboardAction,
		"/" + keyboardActionTypes + "/tab": actions.BuildTabKeyboardAction(),

		"/" + namespaceActionTypes + "/sync":              actions.BuildSyncNamespaceAction(),
		"/" + workload.DeploymentResourceTypes + "/list":  actions.BuildDeploymentListAction(),
		"/" + workload.StatefulSetResourceTypes + "/list": actions.BuildStatefulSetListAction(),
		"/" + workload.DaemonSetResourceTypes + "/list":   actions.BuildDaemonSetListAction(),
		"/" + workload.ReplicaSetResourceTypes + "/list":  actions.BuildReplicaSetListAction(),
		"/" + workload.CronJobResourceTypes + "/list":     actions.BuildCronJobListAction(),
		"/" + workload.JobResourceTypes + "/list":         actions.BuildJobListAction(),
		"/" + service.ServiceResourceTypes + "/list":      actions.BuildServiceListAction(),
		"/" + configMap.ConfigMapResourceTypes + "/list":  actions.BuildConfigMapListAction(),
		"/" + secret.SecretResourceTypes + "/list":        actions.BuildSecretListAction(),
		"/" + pvc.PVCResourceTypes + "/list":              actions.BuildPVCListAction(),
		"/" + pv.PVResourceTypes + "/list":                actions.BuildPVListAction(),
		"/" + node.NodeResourceTypes + "/list":            actions.BuildNodeListAction(),

		"/" + debugMessageActionTypes + "/collect": actions.BuildCollectDebugMessageAction(),
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
			path := el.terminalDashBoard.Menu.Enter()
			el.executeHandler(path)
		case "<Left>":
			path := "/" + keyboardActionTypes + "/left"
			el.executeHandler(path)
		case "<Right>":
			path := "/" + keyboardActionTypes + "/right"
			el.executeHandler(path)
		case "<Up>":
			path := "/" + string(el.getCurrentSelectedPanel()) + "/" + keyboardActionTypes + "/up"
			el.executeHandler(path)
		case "<Down>":
			path := "/" + string(el.getCurrentSelectedPanel()) + "/" + keyboardActionTypes + "/down"
			el.executeHandler(path)
		}
	}
}

func (el *eventListener) executeHandler(path string) {
	if handler, ok := el.pathHandlerMap[path]; ok {
		go handler(el.ctx, el.terminalDashBoard, el.scrapperManagement, el.debugCollector, el.getCurrentNamespace())
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Info, fmt.Sprintf("excute path: %s", path), "eventListener"))
		return
	}
	el.debugCollector.Collect(debug.NewDebugMessage(debug.Warn, fmt.Sprintf("no action match path: %s", path), "eventListener"))
}

func (el *eventListener) getCurrentNamespace() string {
	return el.terminalDashBoard.NamespaceTab.GetCurrentNamespace()
}

func (el *eventListener) getCurrentSelectedPanel() component.PanelTypes {
	return el.terminalDashBoard.GetCurrentPanelTypes()
}
