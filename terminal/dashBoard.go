package terminal

import (
	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper/service"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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

type TerminalDashBoard struct {
	ResourceTable *widgets.Table
	Grid          *ui.Grid
	NamespaceTab  *widgets.TabPane
	YamlPanel     *widgets.Paragraph
	LogPanel      *widgets.Paragraph
	ResourceList  *widgets.List
	Console       *widgets.List
	Menu          *SideMenu
}

func InitDashBoard() *TerminalDashBoard {

	// init workload rTable
	rTable := widgets.NewTable()
	rTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	rTable.Rows = [][]string{{""}}
	rTable.RowSeparator = false
	rTable.TextAlignment = ui.AlignCenter

	// init namespace tab
	nTab := widgets.NewTabPane()

	// init menu
	menu := BuildMenu()

	// init workload tab
	rList := widgets.NewList()
	rList.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
	rList.Rows = make([]string, len(resourceTypes))
	copy(rList.Rows, resourceTypes)

	// yaml panel
	yPanel := widgets.NewParagraph()
	yPanel.Title = "Resource Yaml"

	// log panel
	lPanel := widgets.NewParagraph()
	lPanel.Title = "Log"

	// debug console
	console := widgets.NewList()
	console.Title = "Debug Console"
	console.Rows = []string{"This is kube-viewer console, every debug message will show here."}

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
			// ui.NewCol(3.0/10,
			// 	ui.NewRow(1.0/4, yPanel),
			// 	ui.NewRow(3.0/4, yPanel),
			// ),
		),
	)

	return &TerminalDashBoard{
		Grid:          grid,
		ResourceTable: rTable,
		Menu:          menu,
		NamespaceTab:  nTab,
		YamlPanel:     yPanel,
		Console:       console,
		LogPanel:      lPanel,
	}
}

type menuItem struct {
	displayName string
	actionPath  string
}

func newMenuItem(displayName string, path string) menuItem {
	return menuItem{
		displayName: displayName,
		actionPath:  path,
	}
}

func (mi menuItem) String() string {
	return mi.displayName
}

type SideMenu struct {
	*widgets.Tree
	// items []*widgets.TreeNode
}

func BuildMenu() *SideMenu {
	var nodes []*widgets.TreeNode
	for _, resource := range resourceTypes {
		newNode := &widgets.TreeNode{
			Value: newMenuItem(resource, "/"+resource+"/list"),
		}
		nodes = append(nodes, newNode)
	}

	items := []*widgets.TreeNode{
		{
			Value: newMenuItem("Service", "/"+service.ServiceResourceTypes+"/list"),
			Nodes: nil,
		},
		{
			Value: newMenuItem("Workload", ""),
			Nodes: nodes,
		},
	}

	menu := widgets.NewTree()
	menu.SetNodes(items)
	menu.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)

	return &SideMenu{Tree: menu}
}

// TODO: these functions are too ugly, we will build structure for every component, each component has own function

func (t *TerminalDashBoard) SelectUp() {
	t.Menu.ScrollUp()
}

func (t *TerminalDashBoard) SelectDown() {
	t.Menu.ScrollDown()
}

func (t *TerminalDashBoard) Resize() {
	termWidth, termHeight := ui.TerminalDimensions()
	t.Grid.SetRect(0, 0, termWidth, termHeight)
	ui.Render(t.Grid)
}

func (t *TerminalDashBoard) MoveTabLeft() {
	t.NamespaceTab.FocusLeft()
}

func (t *TerminalDashBoard) MoveTabRight() {
	t.NamespaceTab.FocusRight()
}

func (t *TerminalDashBoard) ShowDebugInfo(message debug.Message) {
	t.Console.Rows = append(t.Console.Rows, message.Format())
	if len(t.Console.Rows) >= 10 {
		t.Console.Rows = t.Console.Rows[1:]
	}
	t.Console.ScrollDown()
}

func (t *TerminalDashBoard) Enter() string {
	t.Menu.ToggleExpand()
	return t.Menu.SelectedNode().Value.(menuItem).actionPath
}
