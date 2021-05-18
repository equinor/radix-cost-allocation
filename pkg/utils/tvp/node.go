package tvp

import (
	"errors"

	"github.com/equinor/radix-cost-allocation/pkg/repository"
	corev1 "k8s.io/api/core/v1"
)

const (
	nodePoolAnnotation = "agentpool"
)

// NewNodeBulkTvpFromNode builds a NodeBulkTvp from the node.
func NewNodeBulkTvpFromNode(node *corev1.Node) (nodeDto repository.NodeBulkTvp, err error) {
	if node == nil {
		err = errors.New("node is nil")
		return
	}

	nodeDto.Name = node.Name
	if nodePool, ok := node.Labels[nodePoolAnnotation]; ok {
		nodeDto.NodePool = nodePool
	}

	return
}
