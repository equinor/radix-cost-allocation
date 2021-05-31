package listers

import (
	"testing"

	v1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestRadixRegistrationLister(t *testing.T) {
	expectedObj1 := &v1.RadixRegistration{}
	expectedObj2 := &v1.RadixRegistration{}
	otherObj := &corev1.Endpoints{}
	store := setupFakeStoreForTest(
		expectedObj1,
		expectedObj2,
		otherObj,
	)
	lister := NewRadixRegistrationLister(store)
	rrs := lister.List()
	assert.Len(t, rrs, 2)
	assert.Equal(t, expectedObj1, rrs[0])
	assert.Equal(t, expectedObj2, rrs[1])
}
