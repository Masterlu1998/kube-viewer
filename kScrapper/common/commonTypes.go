package common

import "time"

const ScrapInterval = time.Second * 1

type KubernetesData interface{}

type ScrapperArgs interface {
	GetNamespaceField() string
}

type ListScrapperArgs struct {
	Namespace string
}

func (l ListScrapperArgs) GetNamespaceField() string {
	return l.Namespace
}

type DetailScrapperArgs struct {
	Namespace string
	Name      string
}

func (d DetailScrapperArgs) GetNamespaceField() string {
	return d.Namespace
}
