package listers

import "k8s.io/client-go/tools/cache"

func setupFakeStoreForTest(listObjects ...any) cache.Store {
	s := &fakeStore{}
	s.ListFunc = func() []any { return listObjects }
	return s
}

type fakeStore struct {
	cache.FakeCustomStore
}

func (f *fakeStore) Bookmark(rv string) {}

func (f *fakeStore) LastStoreSyncResourceVersion() string { return "" }
