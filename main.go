package main

import (
	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/terminal"
)

func main() {
	kScrapper.Start()
	terminal.Run()
}
