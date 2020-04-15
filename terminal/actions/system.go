package actions

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/namespace"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
	"github.com/Masterlu1998/kube-viewer/terminal/path"
)

func BuildUpResourceListKeyboardAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/resourceList/keyboard/up", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.ResourcePanel.ScrollUp()
		tdb.RenderDashboard()
	})
}

func BuildDownResourceListKeyboardAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/resourceList/keyboard/down", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.ResourcePanel.ScrollDown()
		tdb.RenderDashboard()
	})
}

func BuildUpMenuKeyboardAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/menu/keyboard/up", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.Menu.ScrollUp()
		tdb.RenderDashboard()
	})
}

func BuildDownMenuKeyboardAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/menu/keyboard/down", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.Menu.ScrollDown()
		tdb.RenderDashboard()
	})
}

func BuildUpDetailKeyboardAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/detailParagraph/keyboard/up", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.DetailParagraph.ScrollUp()
		tdb.RenderDashboard()
	})
}

func BuildDownDetailKeyboard(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/detailParagraph/keyboard/down", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.DetailParagraph.ScrollDown()
		tdb.RenderDashboard()
	})
}

func BuildLeftKeyboardAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/keyboard/left", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.NamespaceTab.FocusLeft()
		sm.ResetNamespace(tdb.NamespaceTab.GetCurrentNamespace())
		tdb.RenderDashboard()
	})
}

func BuildRightKeyboardAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/keyboard/right", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.NamespaceTab.FocusRight()
		sm.ResetNamespace(tdb.NamespaceTab.GetCurrentNamespace())
		tdb.RenderDashboard()
	})
}

func BuildTabKeyboardAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/keyboard/tab", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.SwitchNextPanel()
		tdb.ResourcePanel.Reset()
		tdb.RenderDashboard()
	})
}

func BuildBackKeyboardAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/keyboard/back", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		tdb.SwitchGrid(component.MainGridTypes)
		sm.StopMainScrapper()
		tdb.DetailParagraph.Clear()
		tdb.RenderDashboard()
	})
}

func BuildCollectDebugMessageAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/debug/collect", func(
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
	})
}

func BuildSyncNamespaceAction(tree *path.TrieTree) {
	tree.RegisterPathWithHandler("/namespace/sync", func(
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
	})
}
