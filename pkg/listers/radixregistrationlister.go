package listers

import (
	v1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	"k8s.io/client-go/tools/cache"
)

type RadixRegistrationLister interface {
	List() []*v1.RadixRegistration
}

type radixRegistrationLister struct {
	store cache.Store
}

func NewRadixRegistrationLister(store cache.Store) RadixRegistrationLister {
	return &radixRegistrationLister{
		store: store,
	}
}

// List returns radixregistrations in the store
func (pl *radixRegistrationLister) List() []*v1.RadixRegistration {
	objs := pl.store.List()
	rrlist := make([]*v1.RadixRegistration, 0, len(objs))

	for _, obj := range objs {
		if rr, ok := obj.(*v1.RadixRegistration); ok {
			rrlist = append(rrlist, rr)
		}
	}

	return rrlist
}
