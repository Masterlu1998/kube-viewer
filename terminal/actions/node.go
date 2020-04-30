package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/node"
	"github.com/Masterlu1998/kube-viewer/terminal/path"
)

var (
	nodeTableHeader   = []string{"name", "status", "roles", "address", "OSImage"}
	nodeTableColWidth = []int{12, 10, 10, 30, 20}
)

func nodeListDataGetter(c common.KubernetesData) ([]string, [][]string, []int, error) {
	nodeInfos, ok := c.([]node.Info)
	if !ok {
		return nil, nil, nil, errors.New("convert to node.Info failed")
	}

	newNodeTableData := make([][]string, 0)
	for _, nodeInfo := range nodeInfos {
		newNodeTableData = append(newNodeTableData, []string{
			nodeInfo.Name,
			nodeInfo.Status,
			nodeInfo.Roles,
			nodeInfo.Address,
			nodeInfo.OSImage,
		})
	}
	return nodeTableHeader, newNodeTableData, nodeTableColWidth, nil
}

func BuildNodeListAction(tree *path.TrieTree) {
	listResourceAction(nodeListDataGetter, tree, node.NodeListScrapperTypes)
}

func BuildNodeDetailAction(tree *path.TrieTree) {
	detailResourceAction(tree, node.NodeDetailScrapperTypes)
}
