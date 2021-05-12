package sync

import (
	"sync"

	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/models"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type nodeSyncer struct {
	lister     listers.NodeLister
	repository repository.Repository
	mu         sync.Mutex
}

// NewNodeSyncer returns an implementation of Syncer that writes node information
// from the NodeLister to the Repository
func NewNodeSyncer(lister listers.NodeLister, repository repository.Repository) Syncer {
	return &nodeSyncer{
		lister:     lister,
		repository: repository,
	}
}

// Sync writes the current list of nodes to the repository
func (s *nodeSyncer) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Info("Started syncing nodes")
	nodes := s.lister.List()
	nodeDtos, err := buildNodeDtoFromNodeList(nodes)
	if err != nil {
		return err
	}

	if err := s.repository.BulkUpsertNodes(nodeDtos); err != nil {
		return err
	}

	log.Info("Finished syncing nodes")
	return nil
}

func buildNodeDtoFromNodeList(nodes []*v1.Node) ([]models.NodeBulkTvp, error) {
	nodeDtos := make([]models.NodeBulkTvp, 0, len(nodes))

	for _, node := range nodes {
		nodeDto, err := models.NodeBulkTvpFromNode(node)
		if err != nil {
			return nil, err
		}
		nodeDtos = append(nodeDtos, nodeDto)
	}

	return nodeDtos, nil
}
