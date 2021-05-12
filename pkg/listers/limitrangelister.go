package listers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type LimitRangeLister interface {
	List() []*corev1.LimitRange
}

type limitRangeLister struct {
	store cache.Store
}

func NewLimitRangeLister(store cache.Store) LimitRangeLister {
	return &limitRangeLister{
		store: store,
	}
}

// List returns limitranges in the store
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
