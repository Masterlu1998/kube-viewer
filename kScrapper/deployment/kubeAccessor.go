package deployment

import (
	"strconv"
	"time"

	"github.com/Masterlu1998/kube-viewer/dataTypes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type kubeAccessor struct {
	kubernetesClient kubernetes.Interface
}

func (ka *kubeAccessor) getDeployments() ([]dataTypes.DeploymentInfo, error) {
	deployments, err := ka.kubernetesClient.AppsV1().Deployments("").List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentInfos []dataTypes.DeploymentInfo
	for _, d := range deployments.Items {
		dInfo := dataTypes.DeploymentInfo{
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
