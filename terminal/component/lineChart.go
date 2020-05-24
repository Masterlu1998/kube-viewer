package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type lineChart struct {
	*widgets.Plot
}

func buildLineChart() *lineChart {
	plot := widgets.NewPlot()
	plot.Title = "braille-mode Line Chart"
	plot.SetRect(0, 0, 50, 15)
	plot.AxesColor = ui.ColorWhite
	plot.LineColors[0] = ui.ColorGreen
	plot.Data = [][]float64{{0, 0, 0, 0}}
	return &lineChart{
		Plot: plot,
	}
}

func (lc *lineChart) RefreshData(data [][]float64) {
	lc.Data = data
}
