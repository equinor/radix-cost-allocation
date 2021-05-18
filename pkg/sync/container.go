package sync

import (
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/tvpbuilder"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

// ContainerSyncJob writes container information to the repository.
// Implements cron.Job interface required by the Cron scheduler
type ContainerSyncJob struct {
	containerTvpBuilder tvpbuilder.ContainerBulkTvpBuilder
	repository          repository.Repository
	sem                 *semaphore.Weighted
}

// NewContainerSyncJob creates a new ContainerSyncJob
func NewContainerSyncJob(containerTvpBuilder tvpbuilder.ContainerBulkTvpBuilder, repository repository.Repository) *ContainerSyncJob {
	return &ContainerSyncJob{
		containerTvpBuilder: containerTvpBuilder,
		repository:          repository,
		sem:                 semaphore.NewWeighted(1),
	}
}

// Sync writes the current list of containers to the repository
func (s *ContainerSyncJob) Sync() error {
	if !s.sem.TryAcquire(1) {
		return NewSyncAlreadyRunningError("container")
	}
	defer s.sem.Release(1)

	log.Info("Start syncing containers")
	containerDtos, err := s.containerTvpBuilder.Build()
	if err != nil {
		return err
	}

	log.Debugf("Writing %v containers to repository", len(containerDtos))
	if err := s.repository.BulkUpsertContainers(containerDtos); err != nil {
		return err
	}

	log.Info("Finished syncing containers")
	return nil
}
