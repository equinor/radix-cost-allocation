package listers

import (
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	dtoutils "github.com/equinor/radix-cost-allocation/pkg/utils/dto"
	v1 "k8s.io/api/core/v1"
)

type nodeDtoMapperFunc func(*v1.Node) repository.NodeBulkDto

var (
	realNodeDtoMapper nodeDtoMapperFunc = dtoutils.MapNodeBulkDtoFromNode
)

// NodeBulkDtoLister defines a method to build a list of NodeBulkDto objects
type NodeBulkDtoLister interface {
	List() []repository.NodeBulkDto
}

type nodeBulkDtoLister struct {
	lister NodeLister
	mapper nodeDtoMapperFunc
}

// NewNodeBulkDtoLister creates a NodeBulkDtoBuilder that builds NodeBulkDto from a NodeLister
func NewNodeBulkDtoLister(lister NodeLister) NodeBulkDtoLister {
	return &nodeBulkDtoLister{
		lister: lister,
		mapper: realNodeDtoMapper,
	}
}

// Build returns a list of NodeBulkDto resources
func (b *nodeBulkDtoLister) List() []repository.NodeBulkDto {
	return buildNodeDtoFromNodeList(b.lister.List(), b.mapper)
}

func buildNodeDtoFromNodeList(nodes []*v1.Node, mapper nodeDtoMapperFunc) []repository.NodeBulkDto {
	nodeDtos := make([]repository.NodeBulkDto, 0, len(nodes))

	for _, node := range nodes {
		if node != nil {
			nodeDto := mapper(node)
			nodeDtos = append(nodeDtos, nodeDto)
		}
	}

	return nodeDtos
}
