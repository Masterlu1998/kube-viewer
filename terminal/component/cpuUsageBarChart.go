package component

import "github.com/gizak/termui/v3/widgets"

type cpuUsageBarChart struct {
	*widgets.BarChart
}

func buildCPUUsageBarChart() *cpuUsageBarChart {
	cpuBarChart := widgets.NewBarChart()
	cpuBarChart.Title = "CPU Usage"
	cpuBarChart.Data = []float64{}
	cpuBarChart.Labels = []string{}
	cpuBarChart.BarWidth = 12
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
	memoryChart.Title = "Memory Usage"
	memoryChart.Data = []float64{}
	memoryChart.BarWidth = 12
	return &memoryUsageBarChart{
		memoryChart,
	}
}

func (m *memoryUsageBarChart) RefreshDataWithLabel(data []float64, label []string) {
	m.Data = data
	m.Labels = label
}
