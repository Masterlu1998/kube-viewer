package main

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/terminal"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	s, err := kScrapper.NewScrapperManagement()
	if err != nil {
		panic(err)
	}
	err = terminal.Run(ctx, cancel, s)
	if err != nil {
		panic(err)
	}
}
