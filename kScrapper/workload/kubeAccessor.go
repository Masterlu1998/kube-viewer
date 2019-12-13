package workload

import (
	"errors"
	"strconv"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type WorkloadInfo struct {
	Name       string
	Namespace  string
	PodsLive   string
	PodsTotal  string
	CreateTime string
	Images     string
}

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
}

func (ka *kubeAccessor) getWorkloads(workloadTypes, namespace string) ([]WorkloadInfo, error) {
	var workloadInfos []WorkloadInfo
	switch workloadTypes {
	case DeploymentResourceTypes:
		deployments, err := ka.kubernetesClient.AppsV1().Deployments(namespace).List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range deployments.Items {
			sInfo := WorkloadInfo{
				Name:       item.Name,
				Namespace:  item.Namespace,
				PodsLive:   strconv.Itoa(int(item.Status.ReadyReplicas)),
				PodsTotal:  strconv.Itoa(int(item.Status.Replicas)),
				CreateTime: item.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
				Images:     item.Spec.Template.Spec.Containers[0].Image,
			}
			workloadInfos = append(workloadInfos, sInfo)
		}
	case StatefulSetResourceTypes:
		statefulSetList, err := ka.kubernetesClient.AppsV1().StatefulSets(namespace).List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range statefulSetList.Items {
			sInfo := WorkloadInfo{
				Name:       item.Name,
				Namespace:  item.Namespace,
				PodsLive:   strconv.Itoa(int(item.Status.ReadyReplicas)),
				PodsTotal:  strconv.Itoa(int(item.Status.Replicas)),
				CreateTime: item.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
				Images:     item.Spec.Template.Spec.Containers[0].Image,
			}
			workloadInfos = append(workloadInfos, sInfo)
		}
	case DaemonSetResourceTypes:
		daemonSetList, err := ka.kubernetesClient.AppsV1().DaemonSets(namespace).List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range daemonSetList.Items {
			sInfo := WorkloadInfo{
				Name:       item.Name,
				Namespace:  item.Namespace,
				PodsLive:   "per node",
				PodsTotal:  "per node",
				CreateTime: item.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
				Images:     item.Spec.Template.Spec.Containers[0].Image,
			}
			workloadInfos = append(workloadInfos, sInfo)
		}
	case ReplicaSetResourceTypes:
		replicaSetList, err := ka.kubernetesClient.AppsV1().ReplicaSets(namespace).List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range replicaSetList.Items {
			sInfo := WorkloadInfo{
				Name:       item.Name,
				Namespace:  item.Namespace,
				PodsLive:   strconv.Itoa(int(item.Status.ReadyReplicas)),
				PodsTotal:  strconv.Itoa(int(item.Status.Replicas)),
				CreateTime: item.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
				Images:     item.Spec.Template.Spec.Containers[0].Image,
			}
			workloadInfos = append(workloadInfos, sInfo)
		}
	case JobResourceTypes:
		replicaSetList, err := ka.kubernetesClient.BatchV1().Jobs(namespace).List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range replicaSetList.Items {
			sInfo := WorkloadInfo{
				Name:       item.Name,
				Namespace:  item.Namespace,
				PodsLive:   "null",
				PodsTotal:  "null",
				CreateTime: item.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
				Images:     item.Spec.Template.Spec.Containers[0].Image,
			}
			workloadInfos = append(workloadInfos, sInfo)
		}
	case CronJobResourceTypes:
		replicaSetList, err := ka.kubernetesClient.BatchV1beta1().CronJobs(namespace).List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range replicaSetList.Items {
			sInfo := WorkloadInfo{
				Name:       item.Name,
				Namespace:  item.Namespace,
				PodsLive:   "null",
				PodsTotal:  "null",
				CreateTime: item.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
				Images:     item.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image,
			}
			workloadInfos = append(workloadInfos, sInfo)
		}
	default:
		return nil, errors.New("invalid ")
	}

	return workloadInfos, nil
}
