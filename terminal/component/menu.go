package component

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type sideMenu struct {
	*widgets.Tree
}

func buildSideMenu() *sideMenu {
	var workLoadItem []*widgets.TreeNode
	for _, resource := range resourceTypes {
		newNode := &widgets.TreeNode{
			Value: newMenuItem(resource, "/"+resource+"/list", resource),
		}
		workLoadItem = append(workLoadItem, newNode)
	}

	items := []*widgets.TreeNode{
		{
			Value: newMenuItem("Overview", "/overview/show", ""),
			Nodes: nil,
		},
		{
			Value: newMenuItem("Cluster", "", ""),
			Nodes: []*widgets.TreeNode{
				{Value: newMenuItem("Persistent Volumes", "/PV/list", "PV")},
				{Value: newMenuItem("Nodes", "/Node/list", "Node")},
			},
		},
		{
			Value: newMenuItem("Workload", "", ""),
			Nodes: workLoadItem,
		},
		{
			Value: newMenuItem("Discovery and Load Balancing", "", ""),
			Nodes: []*widgets.TreeNode{
				{Value: newMenuItem("Service", "/Service/list", "Service")},
			},
		},
		{
			Value: newMenuItem("Config and Storage", "", ""),
			Nodes: []*widgets.TreeNode{
				{Value: newMenuItem("ConfigMaps", "/ConfigMap/list", "ConfigMap")},
				{Value: newMenuItem("Persistent Volume Claims", "/PVC/list", "PVC")},
				{Value: newMenuItem("Secrets", "/Secret/list", "Secret")},
			},
		},
	}

	menu := widgets.NewTree()
	menu.SetNodes(items)
	menu.SelectedRowStyle = termui.NewStyle(termui.ColorYellow)
	menu.ExpandAll()
	menu.Title = "Menu"

	return &sideMenu{Tree: menu}
}

func (m *sideMenu) selectedToggle() {
	if m.BorderStyle == termui.NewStyle(selectedPanelColor) {
		m.BorderStyle = termui.NewStyle(termui.ColorClear)
	} else {
		m.BorderStyle = termui.NewStyle(selectedPanelColor)
	}
}

func (m *sideMenu) Enter() string {
	return m.SelectedNode().Value.(menuItem).actionPath
}

func (m *sideMenu) GetSelectedResourceTypes() string {
	return m.SelectedNode().Value.(menuItem).resourceTypes
}

type menuItem struct {
	displayName   string
	actionPath    string
	resourceTypes string
}

func newMenuItem(displayName string, path string, resourceTypes string) menuItem {
	return menuItem{
		displayName:   displayName,
		actionPath:    path,
		resourceTypes: resourceTypes,
	}
}

func (mi menuItem) String() string {
	return mi.displayName
}
