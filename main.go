package main

import (
	"context"
	"fmt"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/terminal"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	sm, err := kScrapper.NewScrapperManagement()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = terminal.Run(ctx, cancel, sm)
	if err != nil {
		fmt.Println(err)
		return
	}
}
