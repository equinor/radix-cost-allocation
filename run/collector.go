package run

import (
	"context"

	"github.com/equinor/radix-cost-allocation/config"
	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/reflectorcontroller"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/sync"
	"github.com/equinor/radix-cost-allocation/pkg/utils/cronlogger"
	kubeUtils "github.com/equinor/radix-cost-allocation/pkg/utils/kube"
	mssqlUtils "github.com/equinor/radix-cost-allocation/pkg/utils/mssql"
	"github.com/equinor/radix-cost-allocation/pkg/utils/reflectorstore"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type syncerJobWrapper struct {
	syncer sync.Syncer
}

func newSyncerJobWrapper(syncer sync.Syncer) *syncerJobWrapper {
	return &syncerJobWrapper{syncer: syncer}
}

func (w *syncerJobWrapper) Run() {
	if err := w.syncer.Sync(); err != nil {
		handleSyncError(err)
	}
}

func handleSyncError(err error) {
	if errors.Is(err, sync.ErrSyncAlreadyRunning) {
		log.Debug().Msg(err.Error())
	} else {
		log.Error().Msg(err.Error())
	}
}

// InitAndStartCollector starts collecting and writing container and node resources to the database
func InitAndStartCollector(ctx context.Context, sqlConfig config.SQLConfig, cronConfig config.CronSchedule, appNameExcludeList []string) error {
	kubeclient, radixclient, err := kubeUtils.GetKubernetesClients()
	if err != nil {
		return errors.WithMessage(err, "failed to get kubernetes clients")
	}

	db, err := mssqlUtils.OpenSQLServer(sqlConfig.Server, sqlConfig.Database, sqlConfig.User, sqlConfig.Password, sqlConfig.Port)
	if err != nil {
		return errors.WithMessage(err, "failed to init database driver")
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Error().Err(err).Msgf("Failed to close db connection")
		}
	}()
	repo := repository.NewSQLRepository(context.Background(), db, sqlConfig.QueryTimeout)

	// Create reflectors and stores
	podReflector, podStore := reflectorstore.NewPodReflectorAndStore(kubeclient)
	nodeReflector, nodeStore := reflectorstore.NewNodeReflectorAndStore(kubeclient)
	limitrangeReflector, limitrangeStore := reflectorstore.NewLimitRangeReflectorAndStore(kubeclient)
	rrReflector, rrStore := reflectorstore.NewRadixRegistrationReflectorAndStore(radixclient)

	// Create and start reflector controller
	reflectorController := reflectorcontroller.New(podReflector, nodeReflector, limitrangeReflector, rrReflector)
	reflectorController.Start()
	defer reflectorController.Stop()

	// Create listers
	podLister := listers.NewPodLister(podStore)
	nodeLister := listers.NewNodeLister(nodeStore)
	limitRangeLister := listers.NewLimitRangeLister(limitrangeStore)
	rrLister := listers.NewRadixRegistrationLister(rrStore)

	containerDtoLister := listers.NewContainerBulkDtoLister(podLister, rrLister, limitRangeLister)
	nodeDtoLister := listers.NewNodeBulkDtoLister(nodeLister)

	// Create sync jobs
	containerSyncJob := sync.NewContainerSyncJob(containerDtoLister, repo, appNameExcludeList)
	nodeSyncJob := sync.NewNodeSyncJob(nodeDtoLister, repo)

	// Create cron scheduler and add sync jobs
	c := cron.New(cron.WithSeconds(), cron.WithLogger(cronlogger.New(zerolog.Ctx(ctx))))
	if _, err := c.AddJob(cronConfig.PodSync, newSyncerJobWrapper(containerSyncJob)); err != nil {
		return err
	}
	if _, err := c.AddJob(cronConfig.NodeSync, newSyncerJobWrapper(nodeSyncJob)); err != nil {
		return err
	}
	c.Start()
	defer c.Stop()

	<-ctx.Done()
	return nil
}
