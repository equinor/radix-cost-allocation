package sync

import (
	"errors"
	"testing"
	"time"

	mocklisters "github.com/equinor/radix-cost-allocation/pkg/listers/mock"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	mockrepository "github.com/equinor/radix-cost-allocation/pkg/repository/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNodeSyncJob(t *testing.T) {

	t.Run("repository called with correct value", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		values := []repository.NodeBulkDto{
			{Name: "1"},
			{Name: "2"},
		}

		nodeDtoBuilder := mocklisters.NewMockNodeBulkDtoLister(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		nodeDtoBuilder.EXPECT().List().Return(values).Times(1)
		repo.EXPECT().BulkUpsertNodes(values).Return(nil).Times(1)
		job := NewNodeSyncJob(nodeDtoBuilder, repo)
		err := job.Sync()
		assert.Nil(t, err)
	})

	t.Run("SyncAlreadyRunningError returned with second call to Sync", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		nodeDtoBuilder := mocklisters.NewMockNodeBulkDtoLister(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		nodeDtoBuilder.EXPECT().List().Return(nil).Times(1)
		repo.EXPECT().BulkUpsertNodes(nil).DoAndReturn(
			func(arg interface{}) interface{} {
				time.Sleep(100 * time.Millisecond) // Emulate delay in call to repository
				return nil
			},
		).Times(1)
		job := NewNodeSyncJob(nodeDtoBuilder, repo)
		done := make(chan struct{})
		go func() {
			err := job.Sync()
			assert.Nil(t, err)
			close(done)
		}()
		time.Sleep(50 * time.Millisecond)
		err := job.Sync()
		assert.Equal(t, NewSyncAlreadyRunningError("node"), err)
		<-done
	})

	t.Run("repository returns error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		theError := errors.New("an error")
		nodeDtoBuilder := mocklisters.NewMockNodeBulkDtoLister(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		nodeDtoBuilder.EXPECT().List().Return(nil).Times(1)
		repo.EXPECT().BulkUpsertNodes(gomock.Any()).Return(theError).Times(1)
		job := NewNodeSyncJob(nodeDtoBuilder, repo)
		err := job.Sync()
		assert.Equal(t, theError, err)
	})
}
