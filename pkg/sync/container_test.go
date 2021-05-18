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

func TestContainerSyncJob(t *testing.T) {

	t.Run("repository called with correct value", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		values := []repository.ContainerBulkTvp{
			{ContainerID: "1"},
			{ContainerID: "2"},
		}
		defer ctrl.Finish()
		containerTvpBuilder := mocktvpbuilder.NewMockContainerBulkTvpBuilder(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		containerTvpBuilder.EXPECT().Build().Return(values, nil).Times(1)
		repo.EXPECT().BulkUpsertContainers(values).Return(nil).Times(1)
		job := NewContainerSyncJob(containerTvpBuilder, repo)
		err := job.Sync()
		assert.Nil(t, err)
	})

	t.Run("SyncAlreadyRunningError returned with second call to Sync", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()
		containerTvpBuilder := mocktvpbuilder.NewMockContainerBulkTvpBuilder(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		containerTvpBuilder.EXPECT().Build().Return(nil, nil).Times(1)
		repo.EXPECT().BulkUpsertContainers(nil).DoAndReturn(
			func(arg interface{}) interface{} {
				time.Sleep(100 * time.Millisecond) // Emulate delay in call to repository
				return nil
			},
		).Times(1)
		job := NewContainerSyncJob(containerTvpBuilder, repo)
		done := make(chan struct{})
		go func() {
			err := job.Sync()
			assert.Nil(t, err)
			close(done)
		}()
		time.Sleep(50 * time.Millisecond)
		err := job.Sync()
		assert.Equal(t, NewSyncAlreadyRunningError("container"), err)
		<-done
	})

	t.Run("tvp builder returns error", func(t *testing.T) {
		t.Parallel()
		theError := errors.New("an error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		containerTvpBuilder := mocktvpbuilder.NewMockContainerBulkTvpBuilder(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		containerTvpBuilder.EXPECT().Build().Return(nil, theError).Times(1)
		repo.EXPECT().BulkUpsertContainers(gomock.Any()).Times(0)
		job := NewContainerSyncJob(containerTvpBuilder, repo)
		err := job.Sync()
		assert.Equal(t, theError, err)
	})

	t.Run("repository returns error", func(t *testing.T) {
		t.Parallel()
		theError := errors.New("an error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		containerTvpBuilder := mocktvpbuilder.NewMockContainerBulkTvpBuilder(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		containerTvpBuilder.EXPECT().Build().Return(nil, nil).Times(1)
		repo.EXPECT().BulkUpsertContainers(gomock.Any()).Return(theError).Times(1)
		job := NewContainerSyncJob(containerTvpBuilder, repo)
		err := job.Sync()
		assert.Equal(t, theError, err)
	})
}
