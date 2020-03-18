package main

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/debug"
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/terminal"
)

func main() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		ctx, cancel := context.WithCancel(context.Background())
	//
	// 		debugCollector := debug.NewDebugCollector()
	//
	// 		sm, err := kScrapper.NewScrapperManagement(ctx, debugCollector)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		err = terminal.Run(ctx, cancel, sm, debugCollector)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }()
	ctx, cancel := context.WithCancel(context.Background())

	debugCollector := debug.NewDebugCollector()

	sm, err := kScrapper.NewScrapperManagement(ctx, debugCollector)
	if err != nil {
		panic(err)
	}
	err = terminal.Run(ctx, cancel, sm, debugCollector)
	if err != nil {
		panic(err)
	}
}
