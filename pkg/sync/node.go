package sync

import (
	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

// NodeSyncJob writes node information to the repository.
// Implements cron.Job interface required by the Cron scheduler
type NodeSyncJob struct {
	nodeDtoLister listers.NodeBulkDtoLister
	repository    repository.Repository
	sem           *semaphore.Weighted
}

// NewNodeSyncJob creates a new NodeSyncJob
func NewNodeSyncJob(nodeDtoLister listers.NodeBulkDtoLister, repository repository.Repository) *NodeSyncJob {
	return &NodeSyncJob{
		nodeDtoLister: nodeDtoLister,
		repository:    repository,
		sem:           semaphore.NewWeighted(1),
	}
}

// Sync writes the current list of nodes to the repository
func (s *NodeSyncJob) Sync() error {
	if !s.sem.TryAcquire(1) {
		return NewSyncAlreadyRunningError("node")
	}
	defer s.sem.Release(1)

	log.Info().Msg("Start syncing nodes")
	nodeDtos := s.nodeDtoLister.List()

	log.Debug().Msgf("Writing %v nodes to repository", len(nodeDtos))
	if err := s.repository.BulkUpsertNodes(nodeDtos); err != nil {
		return err
	}

	log.Info().Msg("Finished syncing nodes")
	return nil
}
