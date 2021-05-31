package listers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// LimitRangeLister defines a method to list LimitRanges from a k8s cluster
type LimitRangeLister interface {
	// List returns a list of LimitRanges
	List() []*corev1.LimitRange
}

type limitRangeLister struct {
	store cache.Store
}

// NewLimitRangeLister returns a LimitRangeLister that uses a store as source for listing LimitRange resources
func NewLimitRangeLister(store cache.Store) LimitRangeLister {
	return &limitRangeLister{
		store: store,
	}
}

// List returns LimitRanges in the store
func (pl *limitRangeLister) List() []*corev1.LimitRange {
	objs := pl.store.List()
	limitRanges := make([]*corev1.LimitRange, 0, len(objs))

	for _, obj := range objs {
		if limitRange, ok := obj.(*corev1.LimitRange); ok {
			limitRanges = append(limitRanges, limitRange)
		}
	}

	return limitRanges
}
