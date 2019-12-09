package common

import "github.com/gizak/termui/v3/widgets"

type Scrapper interface {
	GetScrapperTypes() string
	ScrapeDataIntoCh() error
	GraphAction(l *widgets.List) error
	GetDataCh() <-chan Data
}

type Data interface{}
