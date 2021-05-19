package dto

import (
	"testing"

	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMapNodeBulkDtoFromNode(t *testing.T) {

	t.Run("node with nodepool label", func(t *testing.T) {
		t.Parallel()
		node := &corev1.Node{ObjectMeta: v1.ObjectMeta{Name: "node", Labels: map[string]string{nodePoolAnnotation: "nodepool"}}}
		expected := repository.NodeBulkDto{Name: "node", NodePool: "nodepool"}
		actual := MapNodeBulkDtoFromNode(node)
		assert.Equal(t, expected, actual)
	})

	t.Run("node without nodepool label", func(t *testing.T) {
		t.Parallel()
		node := &corev1.Node{ObjectMeta: v1.ObjectMeta{Name: "node"}}
		expected := repository.NodeBulkDto{Name: "node"}
		actual := MapNodeBulkDtoFromNode(node)
		assert.Equal(t, expected, actual)
	})
}
