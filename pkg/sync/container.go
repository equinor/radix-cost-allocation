package sync

import (
	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

// ContainerSyncJob writes container information to the repository.
// Implements cron.Job interface required by the Cron scheduler
type ContainerSyncJob struct {
	containerDtoLister listers.ContainerBulkDtoLister
	repository         repository.Repository
	sem                *semaphore.Weighted
}

// NewContainerSyncJob creates a new ContainerSyncJob
func NewContainerSyncJob(containerDtoLister listers.ContainerBulkDtoLister, repository repository.Repository) *ContainerSyncJob {
	return &ContainerSyncJob{
		containerDtoLister: containerDtoLister,
		repository:         repository,
		sem:                semaphore.NewWeighted(1),
	}
}

// Sync writes the current list of containers to the repository
func (s *ContainerSyncJob) Sync() error {
	if !s.sem.TryAcquire(1) {
		return NewSyncAlreadyRunningError("container")
	}
	defer s.sem.Release(1)

	log.Info("Start syncing containers")
	containerDtos := s.containerDtoLister.List()

	log.Debugf("Writing %v containers to repository", len(containerDtos))
	if err := s.repository.BulkUpsertContainers(containerDtos); err != nil {
		return err
	}

	log.Info("Finished syncing containers")
	return nil
}
