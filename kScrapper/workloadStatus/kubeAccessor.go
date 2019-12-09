package workloadStatus

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
}

// func (ka *kubeAccessor) getStatefulSet() ([]) {
//
// }

func (ka *kubeAccessor) getDeployments() ([]deploymentInfo, error) {
	deployments, err := ka.kubernetesClient.AppsV1().Deployments("").List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentInfos []deploymentInfo
	for _, d := range deployments.Items {
		dInfo := deploymentInfo{
			Name:       d.Name,
			Namespace:  d.Namespace,
			PodsLive:   int(d.Status.ReadyReplicas),
			PodsTotal:  int(d.Status.Replicas),
			CreateTime: d.ObjectMeta.CreationTimestamp.Format(time.RFC3339),
			Images:     d.Spec.Template.Spec.Containers[0].Image,
		}
		deploymentInfos = append(deploymentInfos, dInfo)
	}

	return deploymentInfos, nil
}
