package jobs

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

// Run implements the cron.Job interface
func (s *ContainerSyncJob) Run() {
	if err := s.writeToRepository(); err != nil {
		log.Error(err)
	}
}

func (s *ContainerSyncJob) writeToRepository() error {
	if !s.sem.TryAcquire(1) {
		log.Debug("Container sync already running")
		return nil
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
