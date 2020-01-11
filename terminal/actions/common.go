package actions

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
)

type ActionArgs interface{}

type listDataGetter func(common.KubernetesData) (header []string, data [][]string, widths []int, err error)
type detailDataGetter func(common.KubernetesData) (string, error)

type ActionHandler func(
	ctx context.Context,
	tdb *component.TerminalDashBoard,
	sm *kScrapper.ScrapperManagement,
	dc *debug.DebugCollector,
	ns string,
	args ActionArgs,
)

func listResourceAction(getter listDataGetter, scrapperTypes string) ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
		args ActionArgs,
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
			tdb.RenderDashboard()
		}
	}
}

type DetailActionArgs struct {
	Namespace string
	Name      string
}

func detailResourceAction(getter detailDataGetter, scrapperTypes string) ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
		args ActionArgs,
	) {
		err := sm.StartSpecificScrapper(ctx, scrapperTypes, ns)
		if err != nil {
			dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), scrapperTypes))
			return
		}

		for c := range sm.GetSpecificScrapperCh(scrapperTypes) {
			yamlData, err := getter(c)
			if err != nil {
				dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), scrapperTypes+"Getter"))
				continue
			}
			tdb.DetailParagraph.RefreshData(yamlData)
			tdb.RenderDashboard()
		}
	}
}
