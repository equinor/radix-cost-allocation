package listers

import (
	"testing"

	"github.com/equinor/radix-cost-allocation/pkg/listers/mock"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/utils/clock"
	radixv1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestContainerBulkDtoLister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rr := &radixv1.RadixRegistration{ObjectMeta: v1.ObjectMeta{Name: "rr1"}}
	lr := &corev1.LimitRange{ObjectMeta: v1.ObjectMeta{Namespace: "lr1"}}
	pod1 := &corev1.Pod{ObjectMeta: v1.ObjectMeta{Name: "pod1"}}
	pod2 := &corev1.Pod{ObjectMeta: v1.ObjectMeta{Name: "pod2"}}
	containerDto1 := repository.ContainerBulkDto{ContainerID: "c1"}
	containerDto2 := repository.ContainerBulkDto{ContainerID: "c2"}
	containerDto3 := repository.ContainerBulkDto{ContainerID: "c3"}
	mapperReturnVals := map[string][]repository.ContainerBulkDto{
		"pod1": {containerDto1},
		"pod2": {containerDto2, containerDto3},
	}
	var mapperPodReceived []*corev1.Pod
	var mapperRrMapReceived []map[string]*radixv1.RadixRegistration
	var mapperLrMapReceived []map[string]*corev1.LimitRange
	mapperRecorder := func(p *corev1.Pod, m1 map[string]*radixv1.RadixRegistration, m2 map[string]*corev1.LimitRange) {
		mapperPodReceived = append(mapperPodReceived, p)
		mapperRrMapReceived = append(mapperRrMapReceived, m1)
		mapperLrMapReceived = append(mapperLrMapReceived, m2)
	}
	var fakeMapper containerDtoMapperFunc = func(p *corev1.Pod, m1 map[string]*radixv1.RadixRegistration, m2 map[string]*corev1.LimitRange, c clock.Clock) []repository.ContainerBulkDto {
		mapperRecorder(p, m1, m2)
		return mapperReturnVals[p.Name]
	}

	podLister := mock.NewMockPodLister(ctrl)
	podLister.EXPECT().List().Return([]*corev1.Pod{pod1, pod2}).Times(1)
	rrLister := mock.NewMockRadixRegistrationLister(ctrl)
	rrLister.EXPECT().List().Return([]*radixv1.RadixRegistration{rr}).Times(1)
	limitRangeLister := mock.NewMockLimitRangeLister(ctrl)
	limitRangeLister.EXPECT().List().Return([]*corev1.LimitRange{lr}).Times(1)

	lister := containerBulkDtoLister{
		podLister:        podLister,
		rrLister:         rrLister,
		limitRangeLister: limitRangeLister,
		mapper:           fakeMapper,
	}
	containerDtos := lister.List()
	assert.Len(t, containerDtos, 3)
	assert.Equal(t, containerDto1, containerDtos[0])
	assert.Equal(t, containerDto2, containerDtos[1])
	assert.Equal(t, containerDto3, containerDtos[2])
	assert.Len(t, mapperPodReceived, 2)
	assert.Equal(t, pod1, mapperPodReceived[0])
	assert.Equal(t, pod2, mapperPodReceived[1])
	assert.Len(t, mapperRrMapReceived, 2)
	assert.Equal(t, rr, mapperRrMapReceived[0]["rr1"])
	assert.Equal(t, rr, mapperRrMapReceived[1]["rr1"])
	assert.Len(t, mapperLrMapReceived, 2)
	assert.Equal(t, lr, mapperLrMapReceived[0]["lr1"])
	assert.Equal(t, lr, mapperLrMapReceived[1]["lr1"])
}
