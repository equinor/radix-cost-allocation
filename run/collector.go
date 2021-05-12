package run

import (
	"context"
	"log"

	"github.com/equinor/radix-cost-allocation/config"
	"github.com/equinor/radix-cost-allocation/pkg/listers"
	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/sync"
	"github.com/equinor/radix-cost-allocation/pkg/syncscheduler"
	kubeUtils "github.com/equinor/radix-cost-allocation/pkg/utils/kube"
	mssqlUtils "github.com/equinor/radix-cost-allocation/pkg/utils/mssql"
	"github.com/equinor/radix-cost-allocation/pkg/utils/reflector"
	"github.com/equinor/radix-cost-allocation/pkg/utils/reflectorstore"
	"github.com/pkg/errors"
)

func InitAndStartCollector(appConfig config.AppConfig, stopCh <-chan struct{}) error {
	kubeclient, radixclient, err := kubeUtils.GetKubernetesClients()
	if err != nil {
		errors.WithMessage(err, "failed to get kubernetes clients")
	}

	db, err := mssqlUtils.OpenSqlServer(appConfig.SQL.Server, appConfig.SQL.Database, appConfig.SQL.User, appConfig.SQL.Password, appConfig.SQL.Port)
	if err != nil {
		errors.WithMessage(err, "failed to init database driver")
	}
	defer db.Close()
	repo := repository.NewSqlRepository(db, appConfig.SQL.QueryTimeout, context.Background())

	podReflector, podStore := reflectorstore.NewPodReflectorAndStore(kubeclient)
	nodeReflector, nodeStore := reflectorstore.NewNodeReflectorAndStore(kubeclient)
	limitrangeReflector, limitrangeStore := reflectorstore.NewLimitRangeReflectorAndStore(kubeclient)
	rrReflector, rrStore := reflectorstore.NewRadixRegistrationReflectorAndStore(radixclient)

	podLister := listers.NewPodLister(podStore)
	nodeLister := listers.NewNodeLister(nodeStore)
	limitRangeLister := listers.NewLimitRangeLister(limitrangeStore)
	rrLister := listers.NewRadixRegistrationLister(rrStore)
	podSyncer := sync.NewPodSyncer(podLister, rrLister, limitRangeLister, repo)
	nodeSyncer := sync.NewNodeSyncer(nodeLister, repo)

	reflectorController := reflector.NewReflectorController(podReflector, nodeReflector, limitrangeReflector, rrReflector)
	reflectorController.Start()
	defer reflectorController.Stop()

	scheduler, err := syncscheduler.NewWithSyncers(
		syncscheduler.NewSyncSchedulerArg(podSyncer, appConfig.PodSyncSchedule),
		syncscheduler.NewSyncSchedulerArg(nodeSyncer, appConfig.NodeSyncSchedule),
	)
	if err != nil {
		log.Fatalf("failed to create syncscheduler : %v", err)
	}
	scheduler.Start()
	defer scheduler.Stop()

	<-stopCh
	return nil
}
