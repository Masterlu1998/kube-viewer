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

type TerminalDashBoard struct {
	*ui.Grid
	Menu          *sideMenu
	ResourceTable *resourceTable
	NamespaceTab  *namespaceTab
	Console       *debugConsole
}

func InitDashBoard() *TerminalDashBoard {
	// init menu
	menu := buildSideMenu()

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

	return &TerminalDashBoard{
		Grid:          grid,
		ResourceTable: rTable,
		Menu:          menu,
		NamespaceTab:  nTab,
		Console:       console,
	}
}

func (t *TerminalDashBoard) Resize() {
	termWidth, termHeight := ui.TerminalDimensions()
	t.Grid.SetRect(0, 0, termWidth, termHeight)
	ui.Render(t)
}
