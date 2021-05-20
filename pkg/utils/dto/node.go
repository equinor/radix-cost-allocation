package dto

import (
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	corev1 "k8s.io/api/core/v1"
)

const (
	nodePoolAnnotation = "agentpool"
)

// MapNodeBulkDtoFromNode builds a NodeBulkDto from the node.
func MapNodeBulkDtoFromNode(node *corev1.Node) (nodeDto repository.NodeBulkDto) {
	nodeDto.Name = node.Name
	if nodePool, ok := node.Labels[nodePoolAnnotation]; ok {
		nodeDto.NodePool = nodePool
	}

	return
}
