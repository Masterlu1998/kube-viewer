package common

import "time"

const ScrapInterval = time.Second * 1

type KubernetesData interface{}

type ScrapperArgs interface{}

type ListScrapperArgs struct {
	Namespace string
}

type DetailScrapperArgs struct {
	Namespace string
	Name      string
}
