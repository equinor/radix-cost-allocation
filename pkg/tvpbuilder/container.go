package tvpbuilder

import (
	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	tvputils "github.com/equinor/radix-cost-allocation/pkg/utils/tvp"
	v1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	corev1 "k8s.io/api/core/v1"
)

// ContainerBulkTvpBuilder defines a method to build a list of ContainerBulkTvp objects
type ContainerBulkTvpBuilder interface {
	Build() ([]repository.ContainerBulkTvp, error)
}

type containerBulkTvpBuilder struct {
	podLister        listers.PodLister
	rrLister         listers.RadixRegistrationLister
	limitRangeLister listers.LimitRangeLister
}

// NewContainerBulk creates a ContainerBulkTvpBuilder that builds ContainerBulkTvp from the lister parameters
func NewContainerBulk(podLister listers.PodLister,
	rrLister listers.RadixRegistrationLister,
	limitRangeLister listers.LimitRangeLister) ContainerBulkTvpBuilder {
	return &containerBulkTvpBuilder{
		podLister:        podLister,
		rrLister:         rrLister,
		limitRangeLister: limitRangeLister,
	}
}

// Build returns a list of ContainerBulkTvp resources
func (b *containerBulkTvpBuilder) Build() ([]repository.ContainerBulkTvp, error) {
	return buildContainerDtoFromPodList(b.podLister.List(), b.rrLister.List(), b.limitRangeLister.List())
}

func buildContainerDtoFromPodList(pods []*corev1.Pod, rrList []*v1.RadixRegistration, limitRanges []*corev1.LimitRange) ([]repository.ContainerBulkTvp, error) {
	containerDtos := make([]repository.ContainerBulkTvp, 0, len(pods))
	rrMap := buildMapOfRadixRegistrations(rrList)
	limitRangeMap := buildMapOfLimitRanges(limitRanges)

	for _, pod := range pods {
		podContainers, err := tvputils.NewContainerBulkTvpFromPod(pod, rrMap, limitRangeMap)
		if err != nil {
			return nil, err
		}
		containerDtos = append(containerDtos, podContainers...)
	}

	return containerDtos, nil
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
