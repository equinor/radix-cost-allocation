package sync

import (
	"errors"
	"testing"
	"time"

	mocklisters "github.com/equinor/radix-cost-allocation/pkg/listers/mock"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	mockrepository "github.com/equinor/radix-cost-allocation/pkg/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestContainerSyncJob(t *testing.T) {

	t.Run("repository called with correct value", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		values := []repository.ContainerBulkDto{
			{ContainerID: "1"},
			{ContainerID: "2"},
		}

		containerDtoBuilder := mocklisters.NewMockContainerBulkDtoLister(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		containerDtoBuilder.EXPECT().List().Return(values).Times(1)
		repo.EXPECT().BulkUpsertContainers(values).Return(nil).Times(1)
		job := NewContainerSyncJob(containerDtoBuilder, repo, nil)
		err := job.Sync()
		assert.Nil(t, err)
	})

	t.Run("using exclude list - repository called with correct value", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		values := []repository.ContainerBulkDto{
			{ContainerID: "1", ApplicationName: "app1"},
			{ContainerID: "2", ApplicationName: "app2"},
			{ContainerID: "3-1", ApplicationName: "app3"},
			{ContainerID: "3-2", ApplicationName: "APP3"},
			{ContainerID: "4", ApplicationName: "app4"},
		}

		expected := []repository.ContainerBulkDto{
			{ContainerID: "1", ApplicationName: "app1"},
			{ContainerID: "2", ApplicationName: "app2"},
		}

		containerDtoBuilder := mocklisters.NewMockContainerBulkDtoLister(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		containerDtoBuilder.EXPECT().List().Return(values).Times(1)
		repo.EXPECT().BulkUpsertContainers(expected).Return(nil).Times(1)
		job := NewContainerSyncJob(containerDtoBuilder, repo, []string{"app3", "APP4"})
		err := job.Sync()
		assert.Nil(t, err)
	})

	t.Run("SyncAlreadyRunningError returned with second call to Sync", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		containerDtoBuilder := mocklisters.NewMockContainerBulkDtoLister(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		containerDtoBuilder.EXPECT().List().Return(nil).Times(1)
		repo.EXPECT().BulkUpsertContainers(gomock.Any()).DoAndReturn(
			func(arg interface{}) interface{} {
				time.Sleep(100 * time.Millisecond) // Emulate delay in call to repository
				return nil
			},
		).Times(1)
		job := NewContainerSyncJob(containerDtoBuilder, repo, nil)
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

	t.Run("repository returns error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		theError := errors.New("an error")
		containerDtoBuilder := mocklisters.NewMockContainerBulkDtoLister(ctrl)
		repo := mockrepository.NewMockRepository(ctrl)
		containerDtoBuilder.EXPECT().List().Return(nil).Times(1)
		repo.EXPECT().BulkUpsertContainers(gomock.Any()).Return(theError).Times(1)
		job := NewContainerSyncJob(containerDtoBuilder, repo, nil)
		err := job.Sync()
		assert.Equal(t, theError, err)
	})
}
