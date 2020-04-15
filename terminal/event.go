package terminal

import (
	"context"
	"fmt"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	"github.com/Masterlu1998/kube-viewer/terminal/actions"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
	"github.com/Masterlu1998/kube-viewer/terminal/path"
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
	pathHandlerTree    *path.TrieTree
	pathHandlerMap     map[string]path.ActionHandler
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
		pathHandlerTree:    path.BuildTrieTree(),
		pathHandlerMap:     make(map[string]path.ActionHandler),
		debugCollector:     dc,
	}
}

func (el *eventListener) Register() {
	actions.BuildUpResourceListKeyboardAction(el.pathHandlerTree)
	actions.BuildDownResourceListKeyboardAction(el.pathHandlerTree)
	actions.BuildUpMenuKeyboardAction(el.pathHandlerTree)
	actions.BuildDownMenuKeyboardAction(el.pathHandlerTree)
	actions.BuildUpDetailKeyboardAction(el.pathHandlerTree)
	actions.BuildDownDetailKeyboard(el.pathHandlerTree)
	actions.BuildLeftKeyboardAction(el.pathHandlerTree)
	actions.BuildRightKeyboardAction(el.pathHandlerTree)
	actions.BuildTabKeyboardAction(el.pathHandlerTree)
	actions.BuildBackKeyboardAction(el.pathHandlerTree)
	actions.BuildOverviewAction(el.pathHandlerTree)
	actions.BuildSyncNamespaceAction(el.pathHandlerTree)
	actions.BuildCollectDebugMessageAction(el.pathHandlerTree)

	actions.BuildDeploymentListAction(el.pathHandlerTree)
	actions.BuildStatefulSetListAction(el.pathHandlerTree)
	actions.BuildDaemonSetListAction(el.pathHandlerTree)
	actions.BuildReplicaSetListAction(el.pathHandlerTree)
	actions.BuildCronJobListAction(el.pathHandlerTree)
	actions.BuildJobListAction(el.pathHandlerTree)
	actions.BuildServiceListAction(el.pathHandlerTree)
	actions.BuildConfigMapListAction(el.pathHandlerTree)
	actions.BuildSecretListAction(el.pathHandlerTree)
	actions.BuildPVCListAction(el.pathHandlerTree)
	actions.BuildPVListAction(el.pathHandlerTree)
	actions.BuildNodeListAction(el.pathHandlerTree)
	actions.BuildDeploymentDetailAction(el.pathHandlerTree)
	actions.BuildStatefulSetDetailAction(el.pathHandlerTree)
	actions.BuildDaemonSetDetailAction(el.pathHandlerTree)
	actions.BuildReplicaSetDetailAction(el.pathHandlerTree)
	actions.BuildCronJobDetailAction(el.pathHandlerTree)
	actions.BuildJobDetailAction(el.pathHandlerTree)
	actions.BuildServiceDetailAction(el.pathHandlerTree)
	actions.BuildConfigMapDetailAction(el.pathHandlerTree)
	actions.BuildSecretDetailAction(el.pathHandlerTree)
	actions.BuildPVCDetailAction(el.pathHandlerTree)
	actions.BuildPVDetailAction(el.pathHandlerTree)
	actions.BuildNodeDetailAction(el.pathHandlerTree)
}

func (el *eventListener) Listen() error {
	el.executeHandler("/"+debugMessageActionTypes+"/collect", nil)
	el.executeHandler("/"+namespaceActionTypes+"/sync", nil)
	el.executeHandler("/"+overviewActionTypes+"/show", common.ListScrapperArgs{Namespace: el.getCurrentNamespace()})
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
			var p string
			switch el.terminalDashBoard.GetCurrentGrid() {
			case component.OverviewGridTypes:
				fallthrough
			case component.MainGridTypes:
				p = "/" + string(el.getCurrentSelectedPanel()) + "/keyboard/up"
			case component.DetailGridTypes:
				p = "/" + string(component.DetailParagraphPanel) + "/keyboard/up"
			}
			el.executeHandler(p, nil)
		case "<Down>":
			var p string
			switch el.terminalDashBoard.GetCurrentGrid() {
			case component.OverviewGridTypes:
				fallthrough
			case component.MainGridTypes:
				p = "/" + string(el.getCurrentSelectedPanel()) + "/keyboard/down"
			case component.DetailGridTypes:
				p = "/" + string(component.DetailParagraphPanel) + "/keyboard/down"
			}
			el.executeHandler(p, nil)
		}
	}
}

func (el *eventListener) executeHandler(p string, args common.ScrapperArgs) {
	handler := el.pathHandlerTree.GetHandlerWithPath(p)
	if handler != nil {
		go handler(
			el.ctx,
			el.terminalDashBoard,
			el.scrapperManagement,
			el.debugCollector,
			args,
		)
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Info, fmt.Sprintf("excute p: %s", p), "eventListener"))
		return
	}

	if handler, ok := el.pathHandlerMap[p]; ok {
		go handler(
			el.ctx,
			el.terminalDashBoard,
			el.scrapperManagement,
			el.debugCollector,
			args,
		)
		el.debugCollector.Collect(debug.NewDebugMessage(debug.Info, fmt.Sprintf("excute path: %s", p), "eventListener"))
		return
	}
	el.debugCollector.Collect(debug.NewDebugMessage(debug.Warn, fmt.Sprintf("no action match p: %s", p), "eventListener"))
}

func (el *eventListener) getCurrentNamespace() string {
	return el.terminalDashBoard.NamespaceTab.GetCurrentNamespace()
}

func (el *eventListener) getCurrentSelectedPanel() component.PanelTypes {
	return el.terminalDashBoard.GetCurrentPanelTypes()
}
