package syncscheduler

import (
	"github.com/equinor/radix-cost-allocation/pkg/sync"
	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

type syncSchedulerArg struct {
	syncer   sync.Syncer
	schedule string
	sem      *semaphore.Weighted
}

func NewSyncSchedulerArg(syncer sync.Syncer, schedule string) *syncSchedulerArg {
	return &syncSchedulerArg{
		syncer:   syncer,
		schedule: schedule,
		sem:      semaphore.NewWeighted(1),
	}
}

type SyncScheduler struct {
	c *cron.Cron
}

func NewWithSyncers(args ...*syncSchedulerArg) (*SyncScheduler, error) {
	c := cron.New(cron.WithSeconds())

	for _, arg := range args {
		_, err := c.AddFunc(arg.schedule, createJobFunc(arg.syncer, arg.sem))
		if err != nil {
			return nil, err
		}
	}

	return &SyncScheduler{c: c}, nil
}

func createJobFunc(syncer sync.Syncer, sem *semaphore.Weighted) cron.FuncJob {
	return cron.FuncJob(
		func() {
			if !sem.TryAcquire(1) {
				return
			}
			defer sem.Release(1)
			if err := syncer.Sync(); err != nil {
				log.Error(err)
			}
		})
}

func (s *SyncScheduler) Start() {
	s.c.Start()
}

func (s *SyncScheduler) Stop() {
	s.c.Stop()
}
