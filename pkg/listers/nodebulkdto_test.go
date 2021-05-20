package listers

import (
	"testing"

	"github.com/equinor/radix-cost-allocation/pkg/listers/mock"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNodeBulkDtoLister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	node1 := &corev1.Node{ObjectMeta: v1.ObjectMeta{Name: "node1"}}
	node2 := &corev1.Node{ObjectMeta: v1.ObjectMeta{Name: "node2"}}
	nodeDto1 := repository.NodeBulkDto{Name: "node1", NodePool: "pool1"}
	nodeDto2 := repository.NodeBulkDto{Name: "node2", NodePool: "pool2"}
	mapperReturnVals := map[string]repository.NodeBulkDto{
		"node1": nodeDto1,
		"node2": nodeDto2,
	}
	var mapperReceived []*corev1.Node
	mapperRecorder := func(node *corev1.Node) {
		mapperReceived = append(mapperReceived, node)
	}
	var fakeMapper nodeDtoMapperFunc = func(node *corev1.Node) (nodeDto repository.NodeBulkDto) {
		mapperRecorder(node)
		return mapperReturnVals[node.Name]
	}
	nodeLister := mock.NewMockNodeLister(ctrl)
	nodeLister.EXPECT().List().Return([]*corev1.Node{node1, node2}).Times(1)
	lister := nodeBulkDtoLister{lister: nodeLister, mapper: fakeMapper}
	nodeDtos := lister.List()
	assert.Len(t, nodeDtos, 2)
	assert.Equal(t, nodeDto1, nodeDtos[0])
	assert.Equal(t, nodeDto2, nodeDtos[1])
	assert.Len(t, mapperReceived, 2)
	assert.Equal(t, node1, mapperReceived[0])
	assert.Equal(t, node2, mapperReceived[1])
}
