package actions

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
)

func BuildUpResourceListKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
		args ActionArgs,
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
		ns string,
		args ActionArgs,
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
		ns string,
		args ActionArgs,
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
		ns string,
		args ActionArgs,
	) {
		tdb.Menu.ScrollDown()
		tdb.RenderDashboard()
	}
}

func BuildLeftKeyboardAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
		args ActionArgs,
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
		ns string,
		args ActionArgs,
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
		ns string,
		args ActionArgs,
	) {
		tdb.SwitchNextPanel()
		tdb.ResourcePanel.Reset()
		tdb.RenderDashboard()
	}
}

func BuildCollectDebugMessageAction() ActionHandler {
	return func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		ns string,
		args ActionArgs,
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
		ns string,
		args ActionArgs,
	) {
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
			tdb.RenderDashboard()
		}
	}
}
