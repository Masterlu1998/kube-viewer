package path

import (
	"context"
	"fmt"
	"testing"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/terminal/component"
)

func TestTrieTree(t *testing.T) {
	tt := BuildTrieTree()
	tt.RegisterPathWithHandler("/list/deployment", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		fmt.Println("This is path \"/list/deployment\"'s handler!")
	})

	tt.RegisterPathWithHandler("/list/daemonset", func(
		ctx context.Context,
		tdb *component.TerminalDashBoard,
		sm *kScrapper.ScrapperManagement,
		dc *debug.DebugCollector,
		args common.ScrapperArgs,
	) {
		fmt.Println("This is path \"/list/daemonset\"'s handler!")
	})

	haveNode, remainPath := tt.FindLastMatchPathNode("/list/deployment")
	if haveNode.val != "deployment" || len(remainPath) != 0 || haveNode.handler == nil {
		t.Fatal("test failed")
	}

	handler := tt.GetHandlerWithPath("/list/daemonset")
	if handler == nil {
		t.Fatal("test failed")
	}

	handler(nil, nil, nil, nil, nil)
}
