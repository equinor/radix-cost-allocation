package reflectorstore

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func NewNodeReflectorAndStore(client kubernetes.Interface) (*cache.Reflector, cache.Store) {
	store := newStore()
	reflector := newReflector(store, client.CoreV1().RESTClient(), &corev1.Node{}, "nodes", "")
	return reflector, store
}
