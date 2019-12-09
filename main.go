package main

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/terminal"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	s, err := kScrapper.NewScrapperController(ctx)
	if err != nil {
		panic(err)
	}
	err = terminal.Run(cancel, s)
	if err != nil {
		panic(err)
	}
}
