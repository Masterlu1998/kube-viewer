package kScrapper

import (
	"context"
	"time"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workloadStatus"
	"github.com/Masterlu1998/kube-viewer/kube"
	"github.com/sirupsen/logrus"
)

const scrapInterval = time.Second * 1

type ScrapperController struct {
	ScrapperMap map[string]common.Scrapper
}

func NewScrapperController(ctx context.Context) (*ScrapperController, error) {
	s := &ScrapperController{
		ScrapperMap: make(map[string]common.Scrapper),
	}

	kubeClient, err := kube.GetKubernetesClient()
	if err != nil {
		return nil, err
	}
	s.RegisterScrappers(workloadStatus.NewWorkloadStatusScrapper(kubeClient))

	ticker := time.NewTicker(scrapInterval)

	s.StartScrapper(ctx, ticker.C)
	return s, nil
}

func (s *ScrapperController) RegisterScrappers(sps ...common.Scrapper) {
	for _, sp := range sps {
		s.ScrapperMap[sp.GetScrapperTypes()] = sp
	}
}

func (s *ScrapperController) StartScrapper(ctx context.Context, c <-chan time.Time) {
	for _, scrapper := range s.ScrapperMap {
		go func() {
			for {
				select {
				case <-ctx.Done():
					logrus.Infof("stop the %s scrapper", scrapper.GetScrapperTypes())
					return
				case <-c:
					err := scrapper.ScrapeDataIntoCh()
					if err != nil {
						logrus.Errorf("%s scrapper get Data error: %s", err)
					}
				}
			}
		}()
	}
}
