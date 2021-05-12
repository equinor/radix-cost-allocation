package listers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type PodLister interface {
	List() []*corev1.Pod
}

type podLister struct {
	store cache.Store
}

// NewPodLister creates a new PodLister
func NewPodLister(store cache.Store) PodLister {
	return &podLister{
		store: store,
	}
}

// List returns pods in the store
func (pl *podLister) List() []*corev1.Pod {
	objs := pl.store.List()
	pods := make([]*corev1.Pod, 0, len(objs))

	for _, obj := range objs {
		if pod, ok := obj.(*corev1.Pod); ok {
			pods = append(pods, pod)
		}
	}

	return pods
}
