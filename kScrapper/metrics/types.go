package metrics

import "k8s.io/apimachinery/pkg/api/resource"

type NodeMetricsInfo struct {
	Name        string
	CPUUsage    *resource.Quantity
	CPUTotal    *resource.Quantity
	MemoryUsage *resource.Quantity
	MemoryTotal *resource.Quantity
}
