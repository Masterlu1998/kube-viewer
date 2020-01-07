package component

import (
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	ui "github.com/gizak/termui/v3"
)

var (
	resourceTypes = []string{
		workload.DeploymentResourceTypes,
		workload.StatefulSetResourceTypes,
		workload.DaemonSetResourceTypes,
		workload.ReplicaSetResourceTypes,
		workload.CronJobResourceTypes,
		workload.JobResourceTypes,
	}
)

type PanelTypes string

const (
	MenuPanel         PanelTypes = "menu"
	ResourceListPanel PanelTypes = "resourceList"
)

type panelNode struct {
	next  *panelNode
	types PanelTypes
}

var panelIndex = []PanelTypes{MenuPanel, ResourceListPanel}

const selectedPanelColor = ui.ColorYellow

type TerminalDashBoard struct {
	*ui.Grid
	Menu          *sideMenu
	ResourceTable *resourceTable
	NamespaceTab  *namespaceTab
	Console       *debugConsole
	selectedPanel *panelNode
}

func InitDashBoard() *TerminalDashBoard {
	// init menu
	menu := buildSideMenu()
	menu.selectedToggle()

	// init workload rTable
	rTable := BuildResourceTable()

	// init namespace tab
	nTab := buildNamespaceTab()

	// debug console
	console := buildDebugConsole()

	// init layout
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.0/12, nTab),
		ui.NewRow(11.0/12,
			ui.NewCol(2.0/10, ui.NewRow(1, menu)),
			ui.NewCol(8.0/10,
				ui.NewRow(2.0/3, rTable),
				ui.NewRow(1.0/3, console),
			),
		),
	)

	// init selected panel list
	headPanelNode := initPanelLinkedList()

	return &TerminalDashBoard{
		Grid:          grid,
		ResourceTable: rTable,
		Menu:          menu,
		NamespaceTab:  nTab,
		Console:       console,
		selectedPanel: headPanelNode,
	}
}

func initPanelLinkedList() *panelNode {
	head := &panelNode{types: MenuPanel}
	cur := head
	for i := 1; i < len(panelIndex); i++ {
		n := &panelNode{types: panelIndex[i]}
		cur.next = n
		cur = cur.next
	}

	cur.next = head
	return head
}

func (t *TerminalDashBoard) Resize() {
	termWidth, termHeight := ui.TerminalDimensions()
	t.Grid.SetRect(0, 0, termWidth, termHeight)
	ui.Render(t)
}

func (t *TerminalDashBoard) SwitchNextPanel() {
	t.selectPanel()
	t.selectedPanel = t.selectedPanel.next
	t.selectPanel()
}

func (t *TerminalDashBoard) GetCurrentPanelTypes() PanelTypes {
	return t.selectedPanel.types
}

func (t *TerminalDashBoard) selectPanel() {
	switch t.selectedPanel.types {
	case MenuPanel:
		t.Menu.selectedToggle()
	case ResourceListPanel:
		t.ResourceTable.selectedToggle()
	}
}
