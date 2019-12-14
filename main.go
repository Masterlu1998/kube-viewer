package main

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/terminal"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	sm, err := kScrapper.NewScrapperManagement(ctx)
	if err != nil {
		panic(err)
	}
	err = terminal.Run(ctx, cancel, sm)
	if err != nil {
		panic(err)
	}
}
