package resource

import (
	"errors"
	"strconv"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
}

func (ka *kubeAccessor) getDeployments() ([]WorkloadInfo, error) {
	deployments, err := ka.kubernetesClient.AppsV1().Deployments("").List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentInfos []WorkloadInfo
	for _, d := range deployments.Items {
		dInfo := WorkloadInfo{
			Kind:       "deployment",
			Name:       d.Name,
			Namespace:  d.Namespace,
			PodsLive:   strconv.Itoa(int(d.Status.ReadyReplicas)),
			PodsTotal:  strconv.Itoa(int(d.Status.Replicas)),
			CreateTime: d.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
			Images:     d.Spec.Template.Spec.Containers[0].Image,
		}
		deploymentInfos = append(deploymentInfos, dInfo)
	}

	return deploymentInfos, nil
}

func (ka *kubeAccessor) getStatefulSets() ([]WorkloadInfo, error) {
	statefulSetList, err := ka.kubernetesClient.AppsV1().StatefulSets("").List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var statefulSetInfos []WorkloadInfo
	for _, s := range statefulSetList.Items {
		sInfo := WorkloadInfo{
			Kind:       "statefulSet",
			Name:       s.Name,
			Namespace:  s.Namespace,
			PodsLive:   strconv.Itoa(int(s.Status.ReadyReplicas)),
			PodsTotal:  strconv.Itoa(int(s.Status.Replicas)),
			CreateTime: s.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
			Images:     s.Spec.Template.Spec.Containers[0].Image,
		}
		statefulSetInfos = append(statefulSetInfos, sInfo)
	}

	return statefulSetInfos, nil
}

func (ka *kubeAccessor) getWorkloads(workloadTypes ResourceTypes) ([]WorkloadInfo, error) {
	var workloadInfos []WorkloadInfo
	switch workloadTypes {
	case DeploymentResourceTypes:
		deployments, err := ka.kubernetesClient.AppsV1().Deployments("").List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range deployments.Items {
			sInfo := WorkloadInfo{
				Kind:       "statefulSet",
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
		statefulSetList, err := ka.kubernetesClient.AppsV1().StatefulSets("").List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range statefulSetList.Items {
			sInfo := WorkloadInfo{
				Kind:       "statefulSet",
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
		daemonSetList, err := ka.kubernetesClient.AppsV1().DaemonSets("").List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range daemonSetList.Items {
			sInfo := WorkloadInfo{
				Kind:       "statefulSet",
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
		replicaSetList, err := ka.kubernetesClient.AppsV1().ReplicaSets("").List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range replicaSetList.Items {
			sInfo := WorkloadInfo{
				Kind:       "statefulSet",
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
		replicaSetList, err := ka.kubernetesClient.BatchV1().Jobs("").List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range replicaSetList.Items {
			sInfo := WorkloadInfo{
				Kind:       "statefulSet",
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
		replicaSetList, err := ka.kubernetesClient.BatchV1beta1().CronJobs("").List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, item := range replicaSetList.Items {
			sInfo := WorkloadInfo{
				Kind:       "statefulSet",
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

	statefulSetList, err := ka.kubernetesClient.AppsV1().StatefulSets("").List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, s := range statefulSetList.Items {
		sInfo := WorkloadInfo{
			Kind:       "statefulSet",
			Name:       s.Name,
			Namespace:  s.Namespace,
			PodsLive:   strconv.Itoa(int(s.Status.ReadyReplicas)),
			PodsTotal:  strconv.Itoa(int(s.Status.Replicas)),
			CreateTime: s.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
			Images:     s.Spec.Template.Spec.Containers[0].Image,
		}
		workloadInfos = append(workloadInfos, sInfo)
	}

	return workloadInfos, nil
}
