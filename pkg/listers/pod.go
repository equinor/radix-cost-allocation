package listers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// PodLister defines a method to list Pods from a k8s cluster
type PodLister interface {
	// List returns a list of Pods
	List() []*corev1.Pod
}

type podLister struct {
	store cache.Store
}

// NewPodLister returns a PodLister that uses a store as source for listing Pod resources
func NewPodLister(store cache.Store) PodLister {
	return &podLister{
		store: store,
	}
}

// List returns Pods in the store
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
