package reflectorstore

import (
	"time"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

func newStore() cache.Store {
	return cache.NewStore(cache.MetaNamespaceKeyFunc)
}

func newReflector(store cache.Store, client cache.Getter, expectedType interface{}, resource, namespace string) *cache.Reflector {
	return cache.NewReflector(
		cache.NewListWatchFromClient(client, resource, namespace, fields.Everything()),
		expectedType,
		store,
		5*time.Minute,
	)
}
