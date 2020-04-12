package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type resourceDetailPanel struct {
	*widgets.List
}

func buildDetailParagraph() *resourceDetailPanel {
	p := widgets.NewList()
	p.Title = "Detail"
	p.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
	return &resourceDetailPanel{
		List: p,
	}
}

func (r *resourceDetailPanel) RefreshData(data string) {
	temp := make([]byte, 0)
	result := make([]string, 0)
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			result = append(result, string(temp))
			temp = make([]byte, 0)
			continue
		}
		temp = append(temp, data[i])
	}
	r.Rows = result
}

func (r *resourceDetailPanel) Clear() {
	r.Rows = []string{""}
}
