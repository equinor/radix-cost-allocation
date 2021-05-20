package listers

import (
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/utils/clock"
	dtoutils "github.com/equinor/radix-cost-allocation/pkg/utils/dto"
	v1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	corev1 "k8s.io/api/core/v1"
)

type containerDtoMapperFunc func(*corev1.Pod, map[string]*v1.RadixRegistration, map[string]*corev1.LimitRange, clock.Clock) []repository.ContainerBulkDto

var (
	realContainerDtoMapper containerDtoMapperFunc = dtoutils.MapContainerBulkDtoFromPod
)

// ContainerBulkDtoLister defines a method to build a list of ContainerBulkDto objects
type ContainerBulkDtoLister interface {
	List() []repository.ContainerBulkDto
}

type containerBulkDtoLister struct {
	podLister        PodLister
	rrLister         RadixRegistrationLister
	limitRangeLister LimitRangeLister
	mapper           containerDtoMapperFunc
	clock            clock.Clock
}

// NewContainerBulkDtoLister creates a ContainerBulkDtoBuilder that builds ContainerBulkDto from the lister parameters
func NewContainerBulkDtoLister(podLister PodLister, rrLister RadixRegistrationLister, limitRangeLister LimitRangeLister) ContainerBulkDtoLister {
	realClock := &clock.RealClock{}
	return &containerBulkDtoLister{
		podLister:        podLister,
		rrLister:         rrLister,
		limitRangeLister: limitRangeLister,
		mapper:           realContainerDtoMapper,
		clock:            realClock,
	}
}

// List returns a list of ContainerBulkDto resources
func (b *containerBulkDtoLister) List() []repository.ContainerBulkDto {
	return buildContainerDtoFromPodList(b.podLister.List(), b.rrLister.List(), b.limitRangeLister.List(), b.mapper, b.clock)
}

func buildContainerDtoFromPodList(pods []*corev1.Pod, rrList []*v1.RadixRegistration, limitRanges []*corev1.LimitRange, mapper containerDtoMapperFunc, clock clock.Clock) []repository.ContainerBulkDto {
	containerDtos := make([]repository.ContainerBulkDto, 0, len(pods))
	rrMap := buildMapOfRadixRegistrations(rrList)
	limitRangeMap := buildMapOfLimitRanges(limitRanges)

	for _, pod := range pods {
		if pod != nil {
			podContainers := mapper(pod, rrMap, limitRangeMap, clock)
			containerDtos = append(containerDtos, podContainers...)
		}
	}

	return containerDtos
}

func buildMapOfRadixRegistrations(rrList []*v1.RadixRegistration) map[string]*v1.RadixRegistration {
	rrMap := make(map[string]*v1.RadixRegistration, len(rrList))
	for _, rr := range rrList {
		rrMap[rr.Name] = rr
	}
	return rrMap
}

func buildMapOfLimitRanges(limitRanges []*corev1.LimitRange) map[string]*corev1.LimitRange {
	limitRangeMap := make(map[string]*corev1.LimitRange, len(limitRanges))
	for _, limitRange := range limitRanges {
		limitRangeMap[limitRange.Namespace] = limitRange
	}
	return limitRangeMap
}
