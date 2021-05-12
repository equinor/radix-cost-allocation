package reflectorstore

import (
	radixv1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	radixclientset "github.com/equinor/radix-operator/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/cache"
)

func NewRadixRegistrationReflectorAndStore(client radixclientset.Interface) (*cache.Reflector, cache.Store) {
	store := newStore()
	reflector := newReflector(store, client.RadixV1().RESTClient(), &radixv1.RadixRegistration{}, "radixregistrations", "")
	return reflector, store
}
