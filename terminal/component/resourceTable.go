package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type resourceTable struct {
	*widgets.Table
}

func BuildResourceTable() *resourceTable {
	rTable := widgets.NewTable()
	rTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	rTable.Rows = [][]string{{""}}
	rTable.RowSeparator = false
	rTable.TextAlignment = ui.AlignCenter

	return &resourceTable{
		Table: rTable,
	}
}

func (r *resourceTable) RefreshTableData(newData [][]string) {
	r.Rows = newData
}
