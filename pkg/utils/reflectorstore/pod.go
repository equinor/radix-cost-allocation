package reflectorstore

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// NewPodReflectorAndStore creates and returns a new store and a reflector using the store
// The reflector keeps an up to date list of Pod resources in k8s
func NewPodReflectorAndStore(client kubernetes.Interface) (*cache.Reflector, cache.Store) {
	store := newStore()
	reflector := newReflector(store, client.CoreV1().RESTClient(), &corev1.Pod{}, "pods", "")
	return reflector, store
}
