package actions

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
	ui "github.com/gizak/termui/v3"
)

func BuildUpResourceListKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
	) {
		tdb.ResourcePanel.ScrollUp()
		ui.Render(tdb)
	}
}

func BuildDownResourceListKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
	) {
		tdb.ResourcePanel.ScrollDown()
		ui.Render(tdb)
	}
}

func BuildUpMenuKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
	) {
		tdb.Menu.ScrollUp()
		ui.Render(tdb)
	}
}

func BuildDownMenuKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
	) {
		tdb.Menu.ScrollDown()
		ui.Render(tdb)
	}
}

func BuildLeftKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
	) {
		tdb.NamespaceTab.FocusLeft()
		sm.ResetNamespace(tdb.NamespaceTab.GetCurrentNamespace())
		ui.Render(tdb)
	}
}

func BuildRightKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
	) {
		tdb.NamespaceTab.FocusRight()
		sm.ResetNamespace(tdb.NamespaceTab.GetCurrentNamespace())
		ui.Render(tdb)
	}
}

func BuildTabKeyboardAction() ActionHandler {
	return func(ctx context.Context, tdb *component.TerminalDashBoard, sm *kScrapper.ScrapperManagement, dc *debug.DebugCollector, ns string) {
		tdb.SwitchNextPanel()
		tdb.ResourcePanel.Reset()
		ui.Render(tdb)
	}
}

func BuildCollectDebugMessageAction() ActionHandler {
	return func(ctx context.Context, tdb *component.TerminalDashBoard, sm *kScrapper.ScrapperManagement, dc *debug.DebugCollector, ns string) {
		for m := range dc.GetDebugMessageCh() {
			tdb.Console.ShowDebugInfo(m)
			ui.Render(tdb)
		}
	}
}

func BuildSyncNamespaceAction() ActionHandler {
	return func(ctx context.Context, tdb *component.TerminalDashBoard, sm *kScrapper.ScrapperManagement, dc *debug.DebugCollector, ns string) {
		err := sm.StartSpecificScrapper(ctx, namespace.NamespaceScrapperTypes, "")
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
			ui.Render(tdb)
		}
	}
}
