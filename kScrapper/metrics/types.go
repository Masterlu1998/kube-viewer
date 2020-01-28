package metrics

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	nodeMetricsGVP = schema.GroupVersionResource{
		Group:    "metrics.k8s.io",
		Version:  "v1beta1",
		Resource: "nodes",
	}

	podMetricsGVP = schema.GroupVersionResource{
		Group:    "metrics.k8s.io",
		Version:  "v1beta1",
		Resource: "pods",
	}
)

type NodeMetricsInfo struct {
	Name        string
	CPUUsage    *resource.Quantity
	CPUTotal    *resource.Quantity
	MemoryUsage *resource.Quantity
	MemoryTotal *resource.Quantity
}

type PodMetricsInfo struct {
	Name        string
	NameSpace   string
	CPUUsage    *resource.Quantity
	MemoryUsage *resource.Quantity
}
