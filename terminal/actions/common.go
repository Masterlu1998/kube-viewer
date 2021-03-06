package actions

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
	"github.com/Masterlu1998/kube-viewer/terminal/path"
)

type listDataGetter func(common.KubernetesData) (header []string, data [][]string, widths []int, err error)

func listResourceAction(getter listDataGetter, tree *path.TrieTree, scrapperTypes string) {
	tree.RegisterPathWithHandler("/"+scrapperTypes[:len(scrapperTypes)-12]+"/list", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		if tdb.GetCurrentGrid() != component.MainGridTypes {
			tdb.SwitchGrid(component.MainGridTypes)
		}

		sm.StopMainScrapper()
		err := sm.StartSpecificScrapper(ctx, scrapperTypes, args)
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
	})
}

type DetailActionArgs struct {
	Namespace string
	Name      string
}

func detailResourceAction(tree *path.TrieTree, scrapperTypes string) {
	tree.RegisterPathWithHandler("/"+scrapperTypes[:len(scrapperTypes)-14]+"/search", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		if tdb.GetCurrentGrid() != component.DetailGridTypes {
			tdb.SwitchGrid(component.DetailGridTypes)
		}

		err := sm.StartSpecificScrapper(ctx, scrapperTypes, args)
		if err != nil {
			dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), scrapperTypes))
			return
		}

		for c := range sm.GetSpecificScrapperCh(scrapperTypes) {
			yamlData, ok := c.(string)
			if !ok {
				dc.Collect(debug.NewDebugMessage(debug.Error, "convert to string failed", scrapperTypes+"Getter"))
				continue
			}
			tdb.DetailParagraph.RefreshData(yamlData)
			tdb.RenderDashboard()
		}
	})
}
