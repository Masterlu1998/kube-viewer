package terminal

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/dataTypes"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	ui "github.com/gizak/termui/v3"
	"github.com/sirupsen/logrus"
)

// we will show graph in terminal here
func Run(cancel context.CancelFunc, s *kScrapper.ScrapperController) error {
	if err := ui.Init(); err != nil {
		logrus.Errorf("failed to initialize termui: %v", err)
		return err
	}
	defer ui.Close()

	tdb := InitDashBoard()

	go deploymentGraphAction(tdb, s.ScrapperMap[dataTypes.DeploymentScrapperTypes])

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
