package listers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestPodLister(t *testing.T) {
	expectedObj1 := &corev1.Pod{}
	expectedObj2 := &corev1.Pod{}
	otherObj := &corev1.Endpoints{}
	store := setupFakeStoreForTest(
		expectedObj1,
		expectedObj2,
		otherObj,
	)
	lister := NewPodLister(store)
	pods := lister.List()
	assert.Len(t, pods, 2)
	assert.Equal(t, expectedObj1, pods[0])
	assert.Equal(t, expectedObj2, pods[1])
}
