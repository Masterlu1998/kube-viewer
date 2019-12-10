package terminal

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/resource"
	ui "github.com/gizak/termui/v3"
)

// we will show graph in terminal here
func Run(ctx context.Context, cancel context.CancelFunc, s *kScrapper.ScrapperManagement) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	tdb := InitDashBoard()

	go workloadGraphAction(ctx, tdb, s.ResourceScrapper, resource.DeploymentResourceTypes)

	for {
		select {
		case e := <-ui.PollEvents():
			if e.Type == ui.MouseEvent {
				cancel()
				return nil

			}
			if e.Type == ui.KeyboardEvent {
				s.ResourceScrapper.StopResourceScrapper()
				go workloadGraphAction(ctx, tdb, s.ResourceScrapper, resource.StatefulSetResourceTypes)
			}
		}
	}
}
