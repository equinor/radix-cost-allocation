package listers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type NodeLister interface {
	List() []*corev1.Node
}

type nodeLister struct {
	store cache.Store
}

func NewNodeLister(store cache.Store) NodeLister {
	return &nodeLister{
		store: store,
	}
}

// List returns nodes in the store
func (pl *nodeLister) List() []*corev1.Node {
	objs := pl.store.List()
	nodes := make([]*corev1.Node, 0, len(objs))

	for _, obj := range objs {
		if node, ok := obj.(*corev1.Node); ok {
			nodes = append(nodes, node)
		}
	}

	return nodes
}
