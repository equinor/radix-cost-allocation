package sync

import (
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/tvpbuilder"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

// NodeSyncJob writes node information to the repository.
// Implements cron.Job interface required by the Cron scheduler
type NodeSyncJob struct {
	nodeTvpBuilder tvpbuilder.NodeBulkTvpBuilder
	repository     repository.Repository
	sem            *semaphore.Weighted
}

// NewNodeSyncJob creates a new NodeSyncJob
func NewNodeSyncJob(nodeTvpBuilder tvpbuilder.NodeBulkTvpBuilder, repository repository.Repository) *NodeSyncJob {
	return &NodeSyncJob{
		nodeTvpBuilder: nodeTvpBuilder,
		repository:     repository,
		sem:            semaphore.NewWeighted(1),
	}
}

// Sync writes the current list of nodes to the repository
func (s *NodeSyncJob) Sync() error {
	if !s.sem.TryAcquire(1) {
		return NewSyncAlreadyRunningError("node")
	}
	defer s.sem.Release(1)
	log.Info("Start syncing nodes")

	nodeDtos, err := s.nodeTvpBuilder.Build()
	if err != nil {
		return err
	}

	log.Debugf("Writing %v nodes to repository", len(nodeDtos))
	if err := s.repository.BulkUpsertNodes(nodeDtos); err != nil {
		return err
	}

	log.Info("Finished syncing nodes")
	return nil
}
