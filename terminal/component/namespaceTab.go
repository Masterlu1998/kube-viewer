package component

import "github.com/gizak/termui/v3/widgets"

type namespaceTab struct {
	*widgets.TabPane
}

func buildNamespaceTab() *namespaceTab {
	nTab := widgets.NewTabPane()
	nTab.TabNames = []string{""}
	return &namespaceTab{
		TabPane: nTab,
	}
}

func (n *namespaceTab) GetCurrentNamespace() string {
	return n.TabNames[n.ActiveTabIndex]
}

func (n *namespaceTab) RefreshNamespace(newNamespaceData []string) {
	n.TabNames = newNamespaceData
}
