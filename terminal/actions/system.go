package actions

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
)

func BuildUpResourceListKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.ResourcePanel.ScrollUp()
		tdb.RenderDashboard()
	}
}

func BuildDownResourceListKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.ResourcePanel.ScrollDown()
		tdb.RenderDashboard()
	}
}

func BuildUpMenuKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.Menu.ScrollUp()
		tdb.RenderDashboard()
	}
}

func BuildDownMenuKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.Menu.ScrollDown()
		tdb.RenderDashboard()
	}
}

func BuildUpDetailKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.DetailParagraph.ScrollUp()
		tdb.RenderDashboard()
	}
}

func BuildDownDetailKeyboard() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.DetailParagraph.ScrollDown()
		tdb.RenderDashboard()
	}
}

func BuildLeftKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.NamespaceTab.FocusLeft()
		sm.ResetNamespace(tdb.NamespaceTab.GetCurrentNamespace())
		tdb.RenderDashboard()
	}
}

func BuildRightKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.NamespaceTab.FocusRight()
		sm.ResetNamespace(tdb.NamespaceTab.GetCurrentNamespace())
		tdb.RenderDashboard()
	}
}

func BuildTabKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.SwitchNextPanel()
		tdb.ResourcePanel.Reset()
		tdb.RenderDashboard()
	}
}

func BuildBackKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.SwitchGrid(component.MainGrid)
		sm.StopMainScrapper()
		tdb.DetailParagraph.Clear()
		tdb.RenderDashboard()
	}
}

func BuildCollectDebugMessageAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		for m := range dc.GetDebugMessageCh() {
			tdb.Console.ShowDebugInfo(m)
			tdb.RenderDashboard()
		}
	}
}

func BuildSyncNamespaceAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		err := sm.StartSpecificScrapper(ctx, namespace.NamespaceScrapperTypes, args)
		if err != nil {
			dc.Collect(debug.NewDebugMessage(debug.Error, err.Error(), "namespaceAction"))
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case ns := <-sm.GetSpecificScrapperCh(namespace.NamespaceScrapperTypes):
				namespaces := []string{""}
				namespaces = append(namespaces, ns.([]string)...)
				tdb.NamespaceTab.RefreshNamespace(namespaces)
			}
			tdb.RenderDashboard()
		}
	}
}
