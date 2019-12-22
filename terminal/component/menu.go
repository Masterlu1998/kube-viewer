package component

import (
	"github.com/Masterlu1998/kube-viewer/kScrapper/configMap"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pvc"
	"github.com/Masterlu1998/kube-viewer/kScrapper/secret"
	"github.com/Masterlu1998/kube-viewer/kScrapper/service"
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
			Value: newMenuItem(resource, "/"+resource+"/list"),
		}
		workLoadItem = append(workLoadItem, newNode)
	}

	items := []*widgets.TreeNode{
		{
			Value: newMenuItem("Overview", ""),
			Nodes: nil,
		},
		{
			Value: newMenuItem("Cluster", ""),
			Nodes: nil,
		},
		{
			Value: newMenuItem("Workload", ""),
			Nodes: workLoadItem,
		},
		{
			Value: newMenuItem("Discovery and Load Balancing", ""),
			Nodes: []*widgets.TreeNode{
				{
					Value: newMenuItem("Service", "/"+service.ServiceResourceTypes+"/list"),
				},
			},
		},
		{
			Value: newMenuItem("Config and Storage", ""),
			Nodes: []*widgets.TreeNode{
				{
					Value: newMenuItem("ConfigMaps", "/"+configMap.ConfigMapResourceTypes+"/list"),
				},
				{
					Value: newMenuItem("Persistent Volume Claims", "/"+pvc.PVCResourceTypes+"/list"),
				},
				{
					Value: newMenuItem("Secrets", "/"+secret.SecretResourceTypes+"/list"),
				},
			},
		},
	}

	menu := widgets.NewTree()
	menu.SetNodes(items)
	menu.SelectedRowStyle = termui.NewStyle(termui.ColorYellow)
	menu.ExpandAll()

	return &sideMenu{Tree: menu}
}

func (m *sideMenu) Enter() string {
	// t.Menu.ToggleExpand()
	return m.SelectedNode().Value.(menuItem).actionPath
}

type menuItem struct {
	displayName string
	actionPath  string
}

func newMenuItem(displayName string, path string) menuItem {
	return menuItem{
		displayName: displayName,
		actionPath:  path,
	}
}

func (mi menuItem) String() string {
	return mi.displayName
}
