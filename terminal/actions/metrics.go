package actions

import (
	"context"

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
		err := sm.StartSpecificScrapper(ctx, metrics.NodeMetricsListScrapperTypes, args)
		if err != nil {
			dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), metrics.NodeMetricsListScrapperTypes))
			return
		}

		for {
			select {
			case rawNodeMetricList := <-sm.GetSpecificScrapperCh(metrics.NodeMetricsListScrapperTypes):
				nodeMetricList := rawNodeMetricList.([]*metrics.NodeMetricsInfo)
				for _, item := range nodeMetricList {
					memoryUsageInt64 := item.MemoryUsage.MilliValue()
					memoryTotalInt64 := item.MemoryTotal.MilliValue()
					_ = float64(memoryUsageInt64) / float64(memoryTotalInt64)

					cpuUsageInt64 := item.CPUUsage.MilliValue()
					cpuTotalInt64 := item.CPUTotal.MilliValue()
					_ = float64(cpuUsageInt64) / float64(cpuTotalInt64)
				}
			}
		}

	}
}
