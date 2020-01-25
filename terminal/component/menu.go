package component

import (
	"github.com/Masterlu1998/kube-viewer/kScrapper/configMap"
	"github.com/Masterlu1998/kube-viewer/kScrapper/node"
	"github.com/Masterlu1998/kube-viewer/kScrapper/pv"
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
				{Value: newMenuItem("Persistent Volumes", "/"+pv.PVResourceTypes+"/list", pv.PVResourceTypes)},
				{Value: newMenuItem("Nodes", "/"+node.NodeResourceTypes+"/list", node.NodeResourceTypes)},
			},
		},
		{
			Value: newMenuItem("Workload", "", ""),
			Nodes: workLoadItem,
		},
		{
			Value: newMenuItem("Discovery and Load Balancing", "", ""),
			Nodes: []*widgets.TreeNode{
				{
					Value: newMenuItem("Service", "/"+service.ServiceResourceTypes+"/list", service.ServiceResourceTypes),
				},
			},
		},
		{
			Value: newMenuItem("Config and Storage", "", ""),
			Nodes: []*widgets.TreeNode{
				{
					Value: newMenuItem("ConfigMaps", "/"+configMap.ConfigMapResourceTypes+"/list", configMap.ConfigMapResourceTypes),
				},
				{
					Value: newMenuItem("Persistent Volume Claims", "/"+pvc.PVCResourceTypes+"/list", pvc.PVCResourceTypes),
				},
				{
					Value: newMenuItem("Secrets", "/"+secret.SecretResourceTypes+"/list", secret.SecretResourceTypes),
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
