package reflectorstore

import (
	radixv1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	radixclientset "github.com/equinor/radix-operator/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/cache"
)

// NewRadixRegistrationReflectorAndStore creates and returns a new store and a reflector using the store
// The reflector keeps an up to date list of RadixRegistration resources in k8s
func NewRadixRegistrationReflectorAndStore(client radixclientset.Interface) (*cache.Reflector, cache.Store) {
	store := newStore()
	reflector := newReflector(store, client.RadixV1().RESTClient(), &radixv1.RadixRegistration{}, "radixregistrations", "")
	return reflector, store
}
