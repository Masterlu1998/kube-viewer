package actions

import (
	"context"
	"math"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/metrics"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
)

func BuildOverviewAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		sm.StopMainScrapper()
		if tdb.GetCurrentGrid() != component.OverviewGrid {
			tdb.SwitchGrid(component.OverviewGrid)
		}

		err := sm.StartSpecificScrapper(ctx, metrics.NodeMetricsListScrapperTypes, args)
		if err != nil {
			dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), metrics.NodeMetricsListScrapperTypes))
			return
		}

		isOpen := true
		for isOpen {
			select {
			case rawNodeMetricList, ok := <-sm.GetSpecificScrapperCh(metrics.NodeMetricsListScrapperTypes):
				if !ok {
					isOpen = false
					continue
				}
				nodeMetricList := rawNodeMetricList.([]*metrics.NodeMetricsInfo)

				cpuData := make([]float64, 0)
				memoryData := make([]float64, 0)
				labels := make([]string, 0)
				for _, item := range nodeMetricList {
					labels = append(labels, item.Name)

					memoryUsageInt64 := item.MemoryUsage.MilliValue()
					memoryTotalInt64 := item.MemoryTotal.MilliValue()
					memoryUsagePercent := float64(memoryUsageInt64) / float64(memoryTotalInt64)
					memoryData = append(memoryData, math.Trunc(memoryUsagePercent*100))

					cpuUsageInt64 := item.CPUUsage.MilliValue()
					cpuTotalInt64 := item.CPUTotal.MilliValue()
					cpuUsagePercent := float64(cpuUsageInt64) / float64(cpuTotalInt64)
					cpuData = append(cpuData, math.Trunc(cpuUsagePercent*100))
				}
				tdb.MemoryUsageBarChart.RefreshDataWithLabel(memoryData, labels)
				tdb.CPUUsageBarChart.RefreshDataWithLabel(cpuData, labels)
				tdb.RenderDashboard()
			}
		}

	}
}
