package actions

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
	ui "github.com/gizak/termui/v3"
)

type dataGetter func(common.KubernetesData) (header []string, data [][]string, widths []int, err error)

type ActionHandler func(
	ctx context.Context,
	tdb *component.TerminalDashBoard,
	sm *kScrapper.ScrapperManagement,
	dc *debug.DebugCollector,
	ns string,
)

func listResourceAction(getter dataGetter, scrapperTypes string) ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
	) {
		sm.StopMainScrapper()
		err := sm.StartSpecificScrapper(ctx, scrapperTypes, ns)
		if err != nil {
			dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), scrapperTypes))
			return
		}

		for c := range sm.GetSpecificScrapperCh(scrapperTypes) {
			tableHeader, tableData, tableColWidth, err := getter(c)
			if err != nil {
				dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), scrapperTypes+"Getter"))
				continue
			}
			tdb.ResourcePanel.RefreshPanelData(tableHeader, tableData, tableColWidth)
			ui.Render(tdb)
		}
	}
}
