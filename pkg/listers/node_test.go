package listers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/apis/discovery"
)

func TestNodeLister(t *testing.T) {
	expectedObj1 := &corev1.Node{}
	expectedObj2 := &corev1.Node{}
	otherObj := &discovery.EndpointSliceList{}
	store := setupFakeStoreForTest(
		expectedObj1,
		expectedObj2,
		otherObj,
	)
	lister := NewNodeLister(store)
	nodes := lister.List()
	assert.Len(t, nodes, 2)
	assert.Equal(t, expectedObj1, nodes[0])
	assert.Equal(t, expectedObj2, nodes[1])
}
