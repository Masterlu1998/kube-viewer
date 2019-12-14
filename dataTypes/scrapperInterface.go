package dataTypes

import (
	"context"
	"time"
)

const ScrapInterval = time.Second * 1

type KubernetesData interface{}

type Scrapper interface {
	GetScrapperTypes() string
	Watch() <-chan KubernetesData
	StartScrapper(ctx context.Context, namespace string)
	StopResourceScrapper()
}
