package terminal

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workloadStatus"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/sirupsen/logrus"
)

// we will show graph in terminal here
func Run(cancel context.CancelFunc, s *kScrapper.ScrapperController) error {
	if err := ui.Init(); err != nil {
		logrus.Errorf("failed to initialize termui: %v", err)
		return err
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "deployment"
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.SetRect(0, 0, 25, 8)

	ui.Render(l)

	go s.ScrapperMap[workloadStatus.WorkloadStatusTypes].GraphAction(l)

	for {
		select {
		case e := <-ui.PollEvents():
			if e.Type == ui.KeyboardEvent {
				cancel()
				return nil
			}
		}
	}
}
