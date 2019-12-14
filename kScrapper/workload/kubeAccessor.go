package workload

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/Masterlu1998/kube-viewer/kube"
	"k8s.io/apimachinery/pkg/labels"
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
	kubernetesLister *kube.KubeLister
}

func (ka *kubeAccessor) getWorkloads(workloadTypes, namespace string) ([]WorkloadInfo, error) {
	var workloadInfos []WorkloadInfo
	switch workloadTypes {
	case DeploymentResourceTypes:
		deploymentList, err := ka.kubernetesLister.DeploymentLister.Deployments(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range deploymentList {
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
		statefulSetList, err := ka.kubernetesLister.StatefulSetLister.StatefulSets(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range statefulSetList {
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
		daemonSetList, err := ka.kubernetesLister.DaemonSetLister.DaemonSets(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range daemonSetList {
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
		replicaSetList, err := ka.kubernetesLister.ReplicaSetsLister.ReplicaSets(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range replicaSetList {
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
		jobList, err := ka.kubernetesLister.JobLister.Jobs(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range jobList {
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
		cronJobList, err := ka.kubernetesLister.CronJobLister.CronJobs(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range cronJobList {
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

	sort.Slice(workloadInfos, func(i, j int) bool {
		return workloadInfos[i].Name < workloadInfos[j].Name
	})

	return workloadInfos, nil
}
