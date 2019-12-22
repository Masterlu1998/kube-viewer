package kube

import (
	"context"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appv1 "k8s.io/client-go/listers/apps/v1"
	batchv1 "k8s.io/client-go/listers/batch/v1"
	batchv1beta1 "k8s.io/client-go/listers/batch/v1beta1"
	corev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type KubeLister struct {
	DeploymentLister  appv1.DeploymentLister
	DaemonSetLister   appv1.DaemonSetLister
	StatefulSetLister appv1.StatefulSetLister
	ReplicaSetsLister appv1.ReplicaSetLister
	CronJobLister     batchv1beta1.CronJobLister
	JobLister         batchv1.JobLister
	NamespaceLister   corev1.NamespaceLister
	ServiceLister     corev1.ServiceLister
	ConfigMapLister   corev1.ConfigMapLister
	SecretLister      corev1.SecretLister
	PVCLister         corev1.PersistentVolumeClaimLister
}

func NewKubeLister(ctx context.Context, client *kubernetes.Clientset) *KubeLister {
	informer := informers.NewSharedInformerFactory(client, 0)

	deployments := informer.Apps().V1().Deployments()
	daemonSets := informer.Apps().V1().DaemonSets()
	statefulSets := informer.Apps().V1().StatefulSets()
	replicaSets := informer.Apps().V1().ReplicaSets()
	cronJobs := informer.Batch().V1beta1().CronJobs()
	jobs := informer.Batch().V1().Jobs()
	namespaces := informer.Core().V1().Namespaces()
	services := informer.Core().V1().Services()
	configMaps := informer.Core().V1().ConfigMaps()
	secrets := informer.Core().V1().Secrets()
	pvcs := informer.Core().V1().PersistentVolumeClaims()

	go informer.Start(ctx.Done())

	cache.WaitForCacheSync(ctx.Done(),
		deployments.Informer().HasSynced,
		daemonSets.Informer().HasSynced,
		statefulSets.Informer().HasSynced,
		replicaSets.Informer().HasSynced,
		cronJobs.Informer().HasSynced,
		jobs.Informer().HasSynced,
		namespaces.Informer().HasSynced,
		services.Informer().HasSynced,
		configMaps.Informer().HasSynced,
		secrets.Informer().HasSynced,
		pvcs.Informer().HasSynced,
	)

	return &KubeLister{
		DeploymentLister:  deployments.Lister(),
		DaemonSetLister:   daemonSets.Lister(),
		StatefulSetLister: statefulSets.Lister(),
		ReplicaSetsLister: replicaSets.Lister(),
		CronJobLister:     cronJobs.Lister(),
		JobLister:         jobs.Lister(),
		NamespaceLister:   namespaces.Lister(),
		ServiceLister:     services.Lister(),
		ConfigMapLister:   configMaps.Lister(),
		SecretLister:      secrets.Lister(),
		PVCLister:         pvcs.Lister(),
	}
}
