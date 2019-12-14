package kube

import (
	"context"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/listers/apps/v1"
	v12 "k8s.io/client-go/listers/batch/v1"
	"k8s.io/client-go/listers/batch/v1beta1"
	v13 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type KubeLister struct {
	DeploymentLister  v1.DeploymentLister
	DaemonSetLister   v1.DaemonSetLister
	StatefulSetLister v1.StatefulSetLister
	ReplicaSetsLister v1.ReplicaSetLister
	CronJobLister     v1beta1.CronJobLister
	JobLister         v12.JobLister
	NamespaceLister   v13.NamespaceLister
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

	go informer.Start(ctx.Done())

	cache.WaitForCacheSync(ctx.Done(),
		deployments.Informer().HasSynced,
		daemonSets.Informer().HasSynced,
		statefulSets.Informer().HasSynced,
		replicaSets.Informer().HasSynced,
		cronJobs.Informer().HasSynced,
		jobs.Informer().HasSynced,
		namespaces.Informer().HasSynced,
	)

	return &KubeLister{
		DeploymentLister:  deployments.Lister(),
		DaemonSetLister:   daemonSets.Lister(),
		StatefulSetLister: statefulSets.Lister(),
		ReplicaSetsLister: replicaSets.Lister(),
		CronJobLister:     cronJobs.Lister(),
		JobLister:         jobs.Lister(),
		NamespaceLister:   namespaces.Lister(),
	}
}
