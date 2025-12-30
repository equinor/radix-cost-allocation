package listers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/apis/discovery"
)

func TestLimitRangeLister(t *testing.T) {
	expectedObj1 := &corev1.LimitRange{}
	expectedObj2 := &corev1.LimitRange{}
	otherObj := &discovery.EndpointSliceList{}

	store := setupFakeStoreForTest(
		expectedObj1,
		expectedObj2,
		otherObj,
	)
	lister := NewLimitRangeLister(store)
	limitranges := lister.List()
	assert.Len(t, limitranges, 2)
	assert.Equal(t, expectedObj1, limitranges[0])
	assert.Equal(t, expectedObj2, limitranges[1])
}
