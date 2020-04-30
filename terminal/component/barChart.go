package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type cpuUsageBarChart struct {
	*widgets.BarChart
}

func buildCPUUsageBarChart() *cpuUsageBarChart {
	cpuBarChart := widgets.NewBarChart()
	cpuBarChart.Title = "CPU Usage (%)"
	cpuBarChart.Data = []float64{}
	cpuBarChart.Labels = []string{}
	cpuBarChart.BarWidth = 10
	cpuBarChart.BarColors = []ui.Color{ui.ColorWhite, ui.ColorWhite}
	cpuBarChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	cpuBarChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
	return &cpuUsageBarChart{
		cpuBarChart,
	}
}

func (c *cpuUsageBarChart) RefreshDataWithLabel(data []float64, label []string) {
	c.Data = data
	c.Labels = label
}

type memoryUsageBarChart struct {
	*widgets.BarChart
}

func buildMemoryUsageBarChart() *memoryUsageBarChart {
	memoryChart := widgets.NewBarChart()
	memoryChart.Title = "Memory Usage (%)"
	memoryChart.Data = []float64{}
	memoryChart.BarWidth = 10
	memoryChart.BarColors = []ui.Color{ui.ColorWhite, ui.ColorWhite}
	memoryChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	memoryChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
	return &memoryUsageBarChart{
		memoryChart,
	}
}

func (m *memoryUsageBarChart) RefreshDataWithLabel(data []float64, label []string) {
	m.Data = data
	m.Labels = label
}
