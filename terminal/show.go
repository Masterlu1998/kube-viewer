package terminal

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	ui "github.com/gizak/termui/v3"
)

// we will show graph in terminal here
func Run(ctx context.Context, cancel context.CancelFunc, sm *kScrapper.ScrapperManagement) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	tdb := InitDashBoard()

	eventListener := newEventListener(ctx, tdb, cancel, sm)

	eventListener.Register()

	return eventListener.Listen()
}
