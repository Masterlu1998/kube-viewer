package component

import "github.com/gizak/termui/v3/widgets"

type namespaceTab struct {
	*widgets.TabPane
}

func buildNamespaceTab() *namespaceTab {
	return &namespaceTab{
		TabPane: widgets.NewTabPane(),
	}
}

func (n *namespaceTab) RefreshNamespace(newNamespaceData []string) {
	n.TabNames = newNamespaceData
}
