package sync

import (
	"sync"

	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/models"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	v1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type podSyncer struct {
	podLister        listers.PodLister
	rrLister         listers.RadixRegistrationLister
	limitRangeLister listers.LimitRangeLister
	repository       repository.Repository
	mu               sync.Mutex
}

// NewPodSyncer returns an implementation of Syncer that writes pod container information
// from the PodLister to the Repository
func NewPodSyncer(podLister listers.PodLister,
	rrLister listers.RadixRegistrationLister,
	limitRangeLister listers.LimitRangeLister,
	repository repository.Repository) Syncer {
	return &podSyncer{
		podLister:        podLister,
		rrLister:         rrLister,
		limitRangeLister: limitRangeLister,
		repository:       repository,
	}
}

// Sync writes the current list of pod containers to the repository
func (s *podSyncer) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Info("Started syncing pods")
	pods := s.podLister.List()
	rrList := s.rrLister.List()
	limitRanges := s.limitRangeLister.List()
	containerDtos, err := buildContainerDtoFromPodList(pods, rrList, limitRanges)
	if err != nil {
		return err
	}

	if err := s.repository.BulkUpsertContainers(containerDtos); err != nil {
		return err
	}

	log.Info("Finished syncing pods")
	return nil
}

func buildContainerDtoFromPodList(pods []*corev1.Pod, rrList []*v1.RadixRegistration, limitRanges []*corev1.LimitRange) ([]models.ContainerBulkTvp, error) {
	containerDtos := make([]models.ContainerBulkTvp, 0, len(pods))
	rrMap := buildMapOfRadixRegistrations(rrList)
	limitRangeMap := buildMapOfLimitRanges(limitRanges)

	for _, pod := range pods {
		podContainers, err := models.ContainerBulkTvpFromPod(pod, rrMap, limitRangeMap)
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
