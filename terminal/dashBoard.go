package terminal

import (
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	pointer = "->"
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
	ResourceTab   *widgets.List
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

	// init workload tab
	rTab := widgets.NewList()
	rTab.Rows = make([]string, len(resourceTypes))
	copy(rTab.Rows, resourceTypes)

	// yaml panel
	yPanel := widgets.NewParagraph()
	yPanel.Title = "Resource Yaml"

	// log panel
	lPanel := widgets.NewParagraph()
	lPanel.Title = "Log"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.0/12, nTab),
		ui.NewRow(11.0/12,
			ui.NewCol(1.0/10, ui.NewRow(1.0, rTab)),
			ui.NewCol(6.0/10, ui.NewRow(1.0, rTable)),
			ui.NewCol(3.0/10,
				ui.NewRow(1.0/4, yPanel),
				ui.NewRow(3.0/4, yPanel),
			),
		),
	)

	return &TerminalDashBoard{
		Grid:          grid,
		ResourceTable: rTable,
		ResourceTab:   rTab,
		NamespaceTab:  nTab,
		YamlPanel:     yPanel,
		LogPanel:      lPanel,
	}
}

func (t *TerminalDashBoard) RemoveResourcePointer(index int) {
	t.ResourceTab.Rows[index] = resourceTypes[index]
}

func (t *TerminalDashBoard) AddResourcePointer(index int) {
	t.ResourceTab.Rows[index] = pointer + resourceTypes[index]
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
