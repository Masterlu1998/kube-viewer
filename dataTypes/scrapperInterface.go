package dataTypes

type Scrapper interface {
	GetScrapperTypes() string
	ScrapeDataIntoCh() error
	GetDataCh() <-chan Data
}

type Data interface{}
