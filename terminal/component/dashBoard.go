package component

import (
	"sync"

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
	MenuPanel            PanelTypes = "menu"
	ResourceListPanel    PanelTypes = "resourceList"
	DetailParagraphPanel PanelTypes = "detailParagraph"
)

type GridTypes int

const (
	MainGridTypes GridTypes = iota
	DetailGridTypes
	OverviewGridTypes
)

type panelNode struct {
	next  *panelNode
	types PanelTypes
}

var panelIndex = []PanelTypes{MenuPanel, ResourceListPanel}

const selectedPanelColor = ui.ColorYellow

var (
	DetailGrid   *ui.Grid
	OverviewGrid *ui.Grid
	MainGrid     *ui.Grid
)

type TerminalDashBoard struct {
	mutex sync.Mutex
	*ui.Grid
	Menu                *sideMenu
	NamespaceTab        *namespaceTab
	Console             *debugConsole
	ResourcePanel       *resourcePanel
	DetailParagraph     *resourceDetailPanel
	CPUUsageBarChart    *cpuUsageBarChart
	MemoryUsageBarChart *memoryUsageBarChart
	CPUResourceGauge    *cpuResourceGauge
	MemoryResourceGauge *memoryResourceGauge
	LineChart           *lineChart
	selectedPanel       *panelNode
	currentGridTypes    GridTypes
}

func InitDashBoard() *TerminalDashBoard {
	// init menu
	menu := buildSideMenu()
	menu.selectedToggle()

	// init workload rTable
	rTable := BuildResourcePanel()

	// init namespace tab
	nTab := buildNamespaceTab()

	// debug console
	console := buildDebugConsole()

	// detail paragraph
	dParagraph := buildDetailParagraph()

	// overview bar chart
	cpuBarChart := buildCPUUsageBarChart()
	memoryBarChart := buildMemoryUsageBarChart()

	// resource gauge
	cpuResourceGauge := buildCPUResourceGauge()
	memoryResourceGauge := buildMemoryResourceGauge()

	// line chart
	lineChart := buildLineChart()

	// init selected panel list
	headPanelNode := initPanelLinkedList()

	tdb := &TerminalDashBoard{
		ResourcePanel:       rTable,
		Menu:                menu,
		NamespaceTab:        nTab,
		Console:             console,
		DetailParagraph:     dParagraph,
		CPUUsageBarChart:    cpuBarChart,
		MemoryUsageBarChart: memoryBarChart,
		CPUResourceGauge:    cpuResourceGauge,
		MemoryResourceGauge: memoryResourceGauge,
		LineChart:           lineChart,
		selectedPanel:       headPanelNode,
	}

	MainGrid = tdb.buildMainGrid()
	DetailGrid = tdb.buildDetailGrid()
	OverviewGrid = tdb.buildOverviewGrid()

	return tdb
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

func (t *TerminalDashBoard) RenderDashboard() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	ui.Render(t)
}

func (t *TerminalDashBoard) SwitchGrid(types GridTypes) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	switch types {
	case MainGridTypes:
		t.Grid = MainGrid
		t.currentGridTypes = MainGridTypes
	case DetailGridTypes:
		t.Grid = DetailGrid
		t.currentGridTypes = DetailGridTypes
	case OverviewGridTypes:
		t.Grid = OverviewGrid
		t.currentGridTypes = OverviewGridTypes
	}
}

func (t *TerminalDashBoard) buildMainGrid() *ui.Grid {
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.0/12, t.NamespaceTab),
		ui.NewRow(11.0/12,
			ui.NewCol(3.0/14, ui.NewRow(1, t.Menu)),
			ui.NewCol(11.0/14,
				ui.NewRow(2.0/3, t.ResourcePanel),
				ui.NewRow(1.0/3, t.Console),
			),
		),
	)
	return grid
}

func (t *TerminalDashBoard) buildDetailGrid() *ui.Grid {
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1, t.DetailParagraph),
	)
	return grid
}

func (t *TerminalDashBoard) buildOverviewGrid() *ui.Grid {
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.0/12, t.NamespaceTab),
		ui.NewRow(11.0/12,
			ui.NewCol(3.0/14, ui.NewRow(1, t.Menu)),
			ui.NewCol(8.0/14,
				ui.NewRow(2.0/10, t.CPUResourceGauge),
				ui.NewRow(2.0/10, t.MemoryResourceGauge),
				ui.NewRow(6.0/10, t.LineChart),
				// ui.NewRow(3.0/10, t.ResourcePanel),
			),
			ui.NewCol(3.0/14,
				ui.NewRow(1.0/2, t.CPUUsageBarChart),
				ui.NewRow(1.0/2, t.MemoryUsageBarChart),
			),
		),
	)

	return grid
}

func (t *TerminalDashBoard) SwitchNextPanel() {
	t.selectPanelToggle()
	t.selectedPanel = t.selectedPanel.next
	t.selectPanelToggle()
}

func (t *TerminalDashBoard) GetCurrentPanelTypes() PanelTypes {
	return t.selectedPanel.types
}

func (t *TerminalDashBoard) selectPanelToggle() {
	switch t.selectedPanel.types {
	case MenuPanel:
		t.Menu.selectedToggle()
	case ResourceListPanel:
		t.ResourcePanel.selectedToggle()
	}
}

func (t *TerminalDashBoard) GetCurrentGrid() GridTypes {
	return t.currentGridTypes
}
