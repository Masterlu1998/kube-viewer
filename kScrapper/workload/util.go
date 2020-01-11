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

type Info struct {
	Name       string
	Namespace  string
	PodsLive   string
	PodsTotal  string
	CreateTime string
	Images     string
}

func getWorkloads(kubernetesClient kubernetes.Interface, kubernetesLister *kube.KubeLister, workloadTypes, namespace string) ([]Info, error) {
	var workloadInfos []Info
	switch workloadTypes {
	case DeploymentResourceTypes:
		deploymentList, err := kubernetesLister.DeploymentLister.Deployments(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range deploymentList {
			sInfo := Info{
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
		statefulSetList, err := kubernetesLister.StatefulSetLister.StatefulSets(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range statefulSetList {
			sInfo := Info{
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
		daemonSetList, err := kubernetesLister.DaemonSetLister.DaemonSets(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range daemonSetList {
			sInfo := Info{
				Name:       item.Name,
				Namespace:  item.Namespace,
				PodsLive:   "",
				PodsTotal:  "",
				CreateTime: item.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
				Images:     item.Spec.Template.Spec.Containers[0].Image,
			}
			workloadInfos = append(workloadInfos, sInfo)
		}
	case ReplicaSetResourceTypes:
		replicaSetList, err := kubernetesLister.ReplicaSetsLister.ReplicaSets(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range replicaSetList {
			sInfo := Info{
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
		jobList, err := kubernetesLister.JobLister.Jobs(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range jobList {
			sInfo := Info{
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
		cronJobList, err := kubernetesLister.CronJobLister.CronJobs(namespace).List(labels.Everything())
		if err != nil {
			return nil, err
		}

		for _, item := range cronJobList {
			sInfo := Info{
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
