package run

import (
	"context"

	"github.com/equinor/radix-cost-allocation/config"
	"github.com/equinor/radix-cost-allocation/pkg/jobs"
	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/reflectorcontroller"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/tvpbuilder"
	kubeUtils "github.com/equinor/radix-cost-allocation/pkg/utils/kube"
	mssqlUtils "github.com/equinor/radix-cost-allocation/pkg/utils/mssql"
	"github.com/equinor/radix-cost-allocation/pkg/utils/reflectorstore"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

// InitAndStartCollector starts collecting and writing container and node resources to the database
func InitAndStartCollector(sqlConfig config.SQLConfig, cronConfig config.CronSchedule, stopCh <-chan struct{}) error {
	kubeclient, radixclient, err := kubeUtils.GetKubernetesClients()
	if err != nil {
		errors.WithMessage(err, "failed to get kubernetes clients")
	}

	db, err := mssqlUtils.OpenSQLServer(sqlConfig.Server, sqlConfig.Database, sqlConfig.User, sqlConfig.Password, sqlConfig.Port)
	if err != nil {
		errors.WithMessage(err, "failed to init database driver")
	}
	defer db.Close()
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

	// Create TVP (Table Valued Parameter) builders
	containerTvpBuilder := tvpbuilder.NewContainerBulk(podLister, rrLister, limitRangeLister)
	nodeTvpBuilder := tvpbuilder.NewNodeBulk(nodeLister)

	// Create sync jobs
	containerSyncJob := jobs.NewContainerSyncJob(containerTvpBuilder, repo)
	nodeSyncJob := jobs.NewNodeSyncJob(nodeTvpBuilder, repo)

	// Create cron scheduler and add sync jobs
	c := cron.New(cron.WithSeconds())
	if _, err := c.AddJob(cronConfig.PodSync, containerSyncJob); err != nil {
		return err
	}
	if _, err := c.AddJob(cronConfig.NodeSync, nodeSyncJob); err != nil {
		return err
	}
	c.Start()
	defer c.Stop()

	<-stopCh
	return nil
}
