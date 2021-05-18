package sync

import (
	"errors"
	"testing"
	"time"

	"github.com/equinor/radix-cost-allocation/pkg/repository"
	mockrepository "github.com/equinor/radix-cost-allocation/pkg/repository/mock"
	mocktvpbuilder "github.com/equinor/radix-cost-allocation/pkg/tvpbuilder/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNodeSyncJob(t *testing.T) {

	t.Run("repository called with correct value", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		values := []repository.NodeBulkTvp{
			{Name: "1"},
			{Name: "2"},
		}
		defer ctrl.Finish()
		nodeTvpBuilder := mocktvpbuilder.NewMockNodeBulkTvpBuilder(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		nodeTvpBuilder.EXPECT().Build().Return(values, nil).Times(1)
		repo.EXPECT().BulkUpsertNodes(values).Return(nil).Times(1)
		job := NewNodeSyncJob(nodeTvpBuilder, repo)
		err := job.Sync()
		assert.Nil(t, err)
	})

	t.Run("SyncAlreadyRunningError returned with second call to Sync", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()
		nodeTvpBuilder := mocktvpbuilder.NewMockNodeBulkTvpBuilder(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		nodeTvpBuilder.EXPECT().Build().Return(nil, nil).Times(1)
		repo.EXPECT().BulkUpsertNodes(nil).DoAndReturn(
			func(arg interface{}) interface{} {
				time.Sleep(100 * time.Millisecond) // Emulate delay in call to repository
				return nil
			},
		).Times(1)
		job := NewNodeSyncJob(nodeTvpBuilder, repo)
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

	t.Run("tvp builder returns error", func(t *testing.T) {
		t.Parallel()
		theError := errors.New("an error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		nodeTvpBuilder := mocktvpbuilder.NewMockNodeBulkTvpBuilder(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		nodeTvpBuilder.EXPECT().Build().Return(nil, theError).Times(1)
		repo.EXPECT().BulkUpsertNodes(gomock.Any()).Times(0)
		job := NewNodeSyncJob(nodeTvpBuilder, repo)
		err := job.Sync()
		assert.Equal(t, theError, err)
	})

	t.Run("repository returns error", func(t *testing.T) {
		t.Parallel()
		theError := errors.New("an error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		nodeTvpBuilder := mocktvpbuilder.NewMockNodeBulkTvpBuilder(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		nodeTvpBuilder.EXPECT().Build().Return(nil, nil).Times(1)
		repo.EXPECT().BulkUpsertNodes(gomock.Any()).Return(theError).Times(1)
		job := NewNodeSyncJob(nodeTvpBuilder, repo)
		err := job.Sync()
		assert.Equal(t, theError, err)
	})
}
