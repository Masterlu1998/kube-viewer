package terminal

import (
	"context"
	"fmt"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
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
	overviewActionTypes     = "overview"
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
		"/" + keyboardActionTypes + "/left":                                                actions.BuildLeftKeyboardAction(),
		"/" + keyboardActionTypes + "/right":                                               actions.BuildRightKeyboardAction(),
		"/" + keyboardActionTypes + "/back":                                                actions.BuildBackKeyboardAction(),
		"/" + string(component.MenuPanel) + "/" + keyboardActionTypes + "/up":              actions.BuildUpMenuKeyboardAction(),
		"/" + string(component.MenuPanel) + "/" + keyboardActionTypes + "/down":            actions.BuildDownMenuKeyboardAction(),
		"/" + string(component.ResourceListPanel) + "/" + keyboardActionTypes + "/up":      actions.BuildUpResourceListKeyboardAction(),
		"/" + string(component.ResourceListPanel) + "/" + keyboardActionTypes + "/down":    actions.BuildDownResourceListKeyboardAction(),
		"/" + string(component.DetailParagraphPanel) + "/" + keyboardActionTypes + "/up":   actions.BuildUpDetailKeyboardAction(),
		"/" + string(component.DetailParagraphPanel) + "/" + keyboardActionTypes + "/down": actions.BuildDownDetailKeyboard(),
		"/" + keyboardActionTypes + "/tab":                                                 actions.BuildTabKeyboardAction(),

		"/" + overviewActionTypes + "/show": actions.BuildOverviewAction(),

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

		"/" + workload.DeploymentResourceTypes + "/search":  actions.BuildDeploymentDetailAction(),
		"/" + workload.StatefulSetResourceTypes + "/search": actions.BuildStatefulSetDetailAction(),
		"/" + workload.DaemonSetResourceTypes + "/search":   actions.BuildDaemonSetDetailAction(),
		"/" + workload.ReplicaSetResourceTypes + "/search":  actions.BuildReplicaSetDetailAction(),
		"/" + workload.CronJobResourceTypes + "/search":     actions.BuildCronJobDetailAction(),
		"/" + workload.JobResourceTypes + "/search":         actions.BuildJobDetailAction(),
		"/" + service.ServiceResourceTypes + "/search":      actions.BuildServiceDetailAction(),
		"/" + configMap.ConfigMapResourceTypes + "/search":  actions.BuildConfigMapDetailAction(),
		"/" + secret.SecretResourceTypes + "/search":        actions.BuildSecretDetailAction(),
		"/" + pvc.PVCResourceTypes + "/search":              actions.BuildPVCDetailAction(),
		"/" + pv.PVResourceTypes + "/search":                actions.BuildPVDetailAction(),
		"/" + node.NodeResourceTypes + "/search":            actions.BuildNodeDetailAction(),

		"/" + debugMessageActionTypes + "/collect": actions.BuildCollectDebugMessageAction(),
	}
}

func (el *eventListener) Listen() error {
	el.executeHandler("/"+debugMessageActionTypes+"/collect", nil)
	el.executeHandler("/"+namespaceActionTypes+"/sync", nil)
	el.executeHandler("/"+overviewActionTypes+"/show", nil)
	for {
		e := <-ui.PollEvents()
		switch e.ID {
		case "<Resize>":
			el.terminalDashBoard.Resize()
		case "q":
			el.cancelFunc()
			return nil
		case "b":
			path := "/" + keyboardActionTypes + "/back"
			el.executeHandler(path, nil)
		case "<Tab>":
			path := "/" + keyboardActionTypes + "/tab"
			el.executeHandler(path, nil)
		case "<Enter>":
			var (
				path string
				args common.ScrapperArgs
			)
			switch el.terminalDashBoard.GetCurrentPanelTypes() {
			case component.MenuPanel:
				path = el.terminalDashBoard.Menu.Enter()
				args = common.ListScrapperArgs{
					Namespace: el.getCurrentNamespace(),
				}
			case component.ResourceListPanel:
				path = "/" + el.terminalDashBoard.Menu.GetSelectedResourceTypes() + "/search"
				ns, name := el.terminalDashBoard.ResourcePanel.GetSelectedUniqueRowNamespaceAndName()
				args = common.DetailScrapperArgs{
					Namespace: ns,
					Name:      name,
				}
			}
			el.executeHandler(path, args)
		case "<Left>":
			path := "/" + keyboardActionTypes + "/left"
			el.executeHandler(path, nil)
		case "<Right>":
			path := "/" + keyboardActionTypes + "/right"
			el.executeHandler(path, nil)
		case "<Up>":
			var path string
			switch el.terminalDashBoard.GetCurrentGrid() {
			case component.OverviewGrid:
				fallthrough
			case component.MainGrid:
				path = "/" + string(el.getCurrentSelectedPanel()) + "/" + keyboardActionTypes + "/up"
			case component.DetailGrid:
				path = "/" + string(component.DetailParagraphPanel) + "/" + keyboardActionTypes + "/up"
			}
			el.executeHandler(path, nil)
		case "<Down>":
			var path string
			switch el.terminalDashBoard.GetCurrentGrid() {
			case component.OverviewGrid:
				fallthrough
			case component.MainGrid:
				path = "/" + string(el.getCurrentSelectedPanel()) + "/" + keyboardActionTypes + "/down"
			case component.DetailGrid:
				path = "/" + string(component.DetailParagraphPanel) + "/" + keyboardActionTypes + "/down"
			}
			el.executeHandler(path, nil)
		}
	}
}

func (el *eventListener) executeHandler(path string, args common.ScrapperArgs) {
	if handler, ok := el.pathHandlerMap[path]; ok {
		go handler(
			el.ctx,
			el.terminalDashBoard,
			el.scrapperManagement,
			el.debugCollector,
			args,
		)
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
