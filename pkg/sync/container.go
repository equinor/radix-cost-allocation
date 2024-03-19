package sync

import (
	"strings"

	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/utils/slice"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

// ContainerSyncJob writes container information to the repository.
// Implements cron.Job interface required by the Cron scheduler
type ContainerSyncJob struct {
	containerDtoLister listers.ContainerBulkDtoLister
	repository         repository.Repository
	appNameExcludeList []string
	sem                *semaphore.Weighted
}

// NewContainerSyncJob creates a new ContainerSyncJob
func NewContainerSyncJob(containerDtoLister listers.ContainerBulkDtoLister, repository repository.Repository, appNameExcludeList []string) *ContainerSyncJob {
	return &ContainerSyncJob{
		containerDtoLister: containerDtoLister,
		repository:         repository,
		appNameExcludeList: appNameExcludeList,
		sem:                semaphore.NewWeighted(1),
	}
}

// Sync writes the current list of containers to the repository
func (s *ContainerSyncJob) Sync() error {
	if !s.sem.TryAcquire(1) {
		return NewSyncAlreadyRunningError("container")
	}
	defer s.sem.Release(1)

	log.Info().Msg("Start syncing containers")
	containerDtos := s.filterContainerByAppNameExcludeList(s.containerDtoLister.List())

	log.Debug().Msgf("Writing %v containers to repository", len(containerDtos))
	if err := s.repository.BulkUpsertContainers(containerDtos); err != nil {
		return err
	}

	log.Info().Msg("Finished syncing containers")
	return nil
}

func (s *ContainerSyncJob) filterContainerByAppNameExcludeList(containers []repository.ContainerBulkDto) []repository.ContainerBulkDto {
	var idx int
	filtered := make([]repository.ContainerBulkDto, len(containers))
	lowerCaseExcludeList := slice.ToLowerCase(s.appNameExcludeList)

	for _, c := range containers {
		if slice.ContainsString(lowerCaseExcludeList, strings.ToLower(c.ApplicationName)) {
			continue
		}

		filtered[idx] = c
		idx++
	}

	return filtered[:idx]
}
