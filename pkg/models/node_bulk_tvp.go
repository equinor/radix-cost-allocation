package models

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
)

const (
	NodePoolAnnotation = "agentpool"
)

type NodeBulkTvp struct {
	Name     string
	NodePool string
}

func NodeBulkTvpFromNode(node *corev1.Node) (nodeDto NodeBulkTvp, err error) {
	if node == nil {
		err = errors.New("node is nil")
		return
	}

	nodeDto.Name = node.Name
	if nodePool, ok := node.Labels[NodePoolAnnotation]; ok {
		nodeDto.NodePool = nodePool
	}

	return
}
