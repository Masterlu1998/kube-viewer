package terminal

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type TerminalDashBoard struct {
	Table *widgets.Table
	Grid  *ui.Grid
}

var resourceTableHeader = []string{"kind", "name", "namespace", "pods", "create time", "images"}

func InitDashBoard() *TerminalDashBoard {

	// init table
	t := widgets.NewTable()
	t.TextStyle = ui.NewStyle(ui.ColorWhite)
	t.Rows = [][]string{resourceTableHeader}
	t.RowSeparator = false
	t.TextAlignment = ui.AlignCenter

	//
	namespaceTab := widgets.NewTabPane("default", "nginx")

	//
	p := widgets.NewParagraph()
	p.Text = "Resource"

	//
	propertyPanel := widgets.NewParagraph()
	propertyPanel.Title = "Resource property"
	propertyPanel.Text = "status: running"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.0/12, namespaceTab),
		ui.NewRow(11.0/12,
			ui.NewCol(1.0/10, ui.NewRow(1.0, p)),
			ui.NewCol(6.0/10, ui.NewRow(1.0, t)),
			ui.NewCol(3.0/10,
				ui.NewRow(1.0/4, propertyPanel),
				ui.NewRow(3.0/4, propertyPanel),
			),
		),
	)

	ui.Render(grid)

	return &TerminalDashBoard{
		Table: t,
		Grid:  grid,
	}
}
