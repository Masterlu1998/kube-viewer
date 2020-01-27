package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type cpuResourceGauge struct {
	*widgets.Gauge
}

func buildCPUResourceGauge() *cpuResourceGauge {
	gauge := widgets.NewGauge()
	gauge.BarColor = ui.ColorWhite
	gauge.Title = "Total CPU Usage"
	gauge.LabelStyle = ui.NewStyle(ui.ColorWhite)

	return &cpuResourceGauge{
		gauge,
	}
}

func (rg *cpuResourceGauge) RefreshData(percent int) {
	rg.Percent = percent
}

type memoryResourceGauge struct {
	*widgets.Gauge
}

func buildMemoryResourceGauge() *memoryResourceGauge {
	gauge := widgets.NewGauge()
	gauge.BarColor = ui.ColorWhite
	gauge.Title = "Total Memory Usage"
	gauge.LabelStyle = ui.NewStyle(ui.ColorWhite)
	return &memoryResourceGauge{
		gauge,
	}
}

func (rg *memoryResourceGauge) RefreshData(percent int) {
	rg.Percent = percent
}
