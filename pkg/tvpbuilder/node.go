package tvpbuilder

import (
	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	tvputils "github.com/equinor/radix-cost-allocation/pkg/utils/tvp"
	v1 "k8s.io/api/core/v1"
)

// NodeBulkTvpBuilder defines a method to build a list of NodeBulkTvp objects
type NodeBulkTvpBuilder interface {
	Build() ([]repository.NodeBulkTvp, error)
}

type nodeBulkTvpBuilder struct {
	lister listers.NodeLister
}

// NewNodeBulk creates a NodeBulkTvpBuilder that builds NodeBulkTvp from a NodeLister
func NewNodeBulk(lister listers.NodeLister) NodeBulkTvpBuilder {
	return &nodeBulkTvpBuilder{
		lister: lister,
	}
}

// Build returns a list of NodeBulkTvp resources
func (b *nodeBulkTvpBuilder) Build() ([]repository.NodeBulkTvp, error) {
	return buildNodeDtoFromNodeList(b.lister.List())
}

func buildNodeDtoFromNodeList(nodes []*v1.Node) ([]repository.NodeBulkTvp, error) {
	nodeDtos := make([]repository.NodeBulkTvp, 0, len(nodes))

	for _, node := range nodes {
		nodeDto, err := tvputils.NewNodeBulkTvpFromNode(node)
		if err != nil {
			return nil, err
		}
		nodeDtos = append(nodeDtos, nodeDto)
	}

	return nodeDtos, nil
}
