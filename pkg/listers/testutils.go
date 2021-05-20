package listers

import "k8s.io/client-go/tools/cache"

func setupFakeStoreForTest(listObjects ...interface{}) cache.Store {
	return &cache.FakeCustomStore{
		ListFunc: func() []interface{} {
			return listObjects
		},
	}
}
