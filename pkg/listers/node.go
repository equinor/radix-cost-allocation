package listers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// NodeLister defines a method to list Nodes from a k8s cluster
type NodeLister interface {
	// List returns a list of Nodes
	List() []*corev1.Node
}

type nodeLister struct {
	store cache.Store
}

// NewNodeLister returns a NodeLister that uses a store as source for listing Node resources
func NewNodeLister(store cache.Store) NodeLister {
	return &nodeLister{
		store: store,
	}
}

// List returns Nodes in the store
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
