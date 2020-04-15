package path

import (
	"context"
	"strings"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
)

type ActionHandler func(
	ctx context.Context,
	tdb *component.TerminalDashBoard,
	sm *kScrapper.ScrapperManagement,
	dc *debug.DebugCollector,
	args common.ScrapperArgs,
)

type pathNode struct {
	next    map[string]*pathNode
	val     string
	handler ActionHandler
}

type TrieTree struct {
	rootNode *pathNode
}

func BuildTrieTree() *TrieTree {
	root := &pathNode{
		next:    nil,
		val:     "",
		handler: nil,
	}
	return &TrieTree{
		rootNode: root,
	}
}

func (tt *TrieTree) RegisterPathWithHandler(path string, handler ActionHandler) {
	cur, remainPathSlice := tt.FindLastMatchPathNode(path)
	if cur == nil {
		cur = tt.rootNode
	}

	for len(remainPathSlice) != 1 {
		section := remainPathSlice[0]
		nextNode := &pathNode{
			next: make(map[string]*pathNode),
			val:  section,
		}

		cur.next[section] = nextNode
		cur = nextNode
		remainPathSlice = remainPathSlice[1:]
	}

	section := remainPathSlice[0]
	cur.next[section] = &pathNode{
		val:     section,
		handler: handler,
	}
}

func (tt *TrieTree) FindLastMatchPathNode(path string) (*pathNode, []string) {
	pathSlice := strings.Split(path, "/")
	pathSlice = pathSlice[1:]
	if len(pathSlice) == 0 {
		return nil, pathSlice
	}

	cur := tt.rootNode
	for len(pathSlice) != 0 {
		if cur.next == nil {
			cur.next = make(map[string]*pathNode)
			break
		}

		nextNode, ok := cur.next[pathSlice[0]]
		if !ok {
			break
		}

		cur = nextNode
		pathSlice = pathSlice[1:]
	}

	return cur, pathSlice
}

func (tt *TrieTree) GetHandlerWithPath(path string) ActionHandler {
	lastMatchNode, remainPathSlice := tt.FindLastMatchPathNode(path)
	if lastMatchNode != nil && len(remainPathSlice) == 0 {
		return lastMatchNode.handler
	}

	return nil
}
