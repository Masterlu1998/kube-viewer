package actions

import (
	"context"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/metrics"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
	"github.com/Masterlu1998/kube-viewer/terminal/path"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	podTableHeader   = []string{"name", "namespace", "cpu usage", "memory usage"}
	podTableColWidth = []int{22, 12, 12, 12}
	CPUUsageCache    = make([]float64, 0, 60)
)

func BuildOverviewAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/overview/show", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		sm.StopMainScrapper()
		if tdb.GetCurrentGrid() != component.OverviewGridTypes {
			tdb.SwitchGrid(component.OverviewGridTypes)
		}

		err := sm.StartSpecificScrapper(ctx, metrics.NodeMetricsListScrapperTypes, args)
		if err != nil {
			dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), metrics.NodeMetricsListScrapperTypes))
			return
		}

		err = sm.StartSpecificScrapper(ctx, metrics.PodMetricsListScrapperTypes, args)
		if err != nil {
			dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), metrics.PodMetricsListScrapperTypes))
			return
		}

		isOpen := true
		for isOpen {
			select {
			case rawPodMetricList, ok := <-sm.GetSpecificScrapperCh(metrics.PodMetricsListScrapperTypes):
				if !ok {
					isOpen = false
					continue
				}

				podMetricList := rawPodMetricList.([]*metrics.PodMetricsInfo)
				data := make([][]string, 0)
				for _, item := range podMetricList {
					col := []string{
						item.Name,
						item.NameSpace,
						strconv.FormatFloat(float64(int(item.CPUUsage.ScaledValue(resource.Micro)))/10000, 'f', 2, 64) + "%",
						strconv.Itoa(int(item.MemoryUsage.Value())/1024) + "MB",
					}
					data = append(data, col)
				}
				tdb.ResourcePanel.RefreshPanelData(podTableHeader, data, podTableColWidth)
			case rawNodeMetricList, ok := <-sm.GetSpecificScrapperCh(metrics.NodeMetricsListScrapperTypes):
				if !ok {
					isOpen = false
					continue
				}
				nodeMetricList := rawNodeMetricList.([]*metrics.NodeMetricsInfo)

				cpuData := make([]float64, 0)
				memoryData := make([]float64, 0)
				var totalCPU, usageCPU, totalMemory, usageMemory int64
				labels := make([]string, 0)
				for _, item := range nodeMetricList {
					labels = append(labels, item.CPUTotal.String()+" Cores")

					memoryUsageInt64 := item.MemoryUsage.MilliValue()
					memoryTotalInt64 := item.MemoryTotal.MilliValue()
					memoryUsagePercent := float64(memoryUsageInt64) / float64(memoryTotalInt64)
					memoryData = append(memoryData, math.Trunc(memoryUsagePercent*100))
					totalMemory += memoryTotalInt64
					usageMemory += memoryUsageInt64

					cpuUsageInt64 := item.CPUUsage.MilliValue()
					cpuTotalInt64 := item.CPUTotal.MilliValue()
					cpuUsagePercent := float64(cpuUsageInt64) / float64(cpuTotalInt64)
					cpuData = append(cpuData, math.Trunc(cpuUsagePercent*100))
					totalCPU += cpuTotalInt64
					usageCPU += cpuUsageInt64
				}
				tdb.MemoryUsageBarChart.RefreshDataWithLabel(memoryData, labels)
				tdb.CPUUsageBarChart.RefreshDataWithLabel(cpuData, labels)
				finalUsageMemoryData := generateIntFloatData(int(usageMemory) * 100 / int(totalMemory))
				finalUsageCPUData := generateIntFloatData((int(usageCPU)) * 100 / int(totalCPU))
				tdb.MemoryResourceGauge.RefreshData(finalUsageMemoryData)
				tdb.CPUResourceGauge.RefreshData(finalUsageCPUData)

				if len(CPUUsageCache) == 60 {
					CPUUsageCache = CPUUsageCache[1:]
				}
				CPUUsageCache = append(CPUUsageCache, float64(finalUsageCPUData))

				if len(CPUUsageCache) >= 2 {
					tdb.LineChart.RefreshData([][]float64{CPUUsageCache})
				}

				tdb.RenderDashboard()
			}
		}
	})
}

func generateIntFloatData(realData int) int {
	rand.Seed(time.Now().UnixNano())
	floatData := rand.Intn(4)

	if rand.Int()%2 == 0 {
		return realData + floatData
	}

	return realData - floatData
}
