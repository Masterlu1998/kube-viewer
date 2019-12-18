package component

import (
	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/gizak/termui/v3/widgets"
)

type debugConsole struct {
	*widgets.List
}

func buildDebugConsole() *debugConsole {
	console := widgets.NewList()
	console.Title = "Debug Console"
	console.Rows = []string{"This is kube-viewer console, every debug message will show here."}
	return &debugConsole{
		List: console,
	}
}

func (d *debugConsole) ShowDebugInfo(message debug.Message) {
	d.Rows = append(d.Rows, message.Format())
	if len(d.Rows) >= 10 {
		d.Rows = d.Rows[1:]
	}
	d.ScrollDown()
}
