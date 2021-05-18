package run

import (
	"time"

	"github.com/equinor/radix-cost-allocation/clients"
	"github.com/equinor/radix-cost-allocation/config"
	"github.com/equinor/radix-cost-allocation/models"
	"github.com/pkg/errors"
	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

var sem = semaphore.NewWeighted(1)

// InitAndStartOldDataCollector starts the old (and soon to be deprecated) resource collector
func InitAndStartOldDataCollector(prometheusAPIURL, cronSchedule string, sqlConfig config.SQLConfig, stopCh <-chan struct{}) {
	promClient := clients.PrometheusClient{Address: prometheusAPIURL}
	sqlClient, err := clients.NewSQLClient(sqlConfig.Server, sqlConfig.Database, sqlConfig.Port, sqlConfig.User, sqlConfig.Password, sqlConfig.QueryTimeout)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sqlClient.Close()

	log.Infof("Registering old job using cron schedule %s", cronSchedule)
	c := cron.New(cron.WithSeconds())
	c.AddFunc(cronSchedule, func() {
		if !sem.TryAcquire(1) {
			return
		}
		defer sem.Release(1)

		if err := moveResourceRequestsFromPrometheusToSQLDB(promClient, sqlClient); err != nil {
			log.Error(err)
		}
	})

	c.Start()
	defer c.Stop()
	<-stopCh
}

func moveResourceRequestsFromPrometheusToSQLDB(promClient clients.PrometheusClient, sqlClient clients.SQLClient) error {
	measuredTimeUTC := time.Now().UTC()
	reqResources, err := promClient.GetRequiredResources(measuredTimeUTC)
	if err != nil {
		return errors.WithMessage(err, "error getting required resources")
	}

	clusterCPUCores, err := promClient.GetClusterTotalCPUCoresFromPrometheus(measuredTimeUTC)
	if err != nil {
		return errors.WithMessage(err, "error getting node cpu count")
	}
	clusterCPUMillieCores := clusterCPUCores * 1000

	clusterMemoryBytes, err := promClient.GetClusterTotalMemoryBytesFromPrometheus(measuredTimeUTC)
	if err != nil {
		return errors.WithMessage(err, "error getting node cpu count")
	}
	clusterMemoryMegaByte := clusterMemoryBytes / 1000000

	runID, err := sqlClient.SaveRun(measuredTimeUTC, clusterCPUMillieCores, clusterMemoryMegaByte)
	if err != nil {
		return errors.WithMessage(err, "error creating Run")
	}
	log.Infof("Run %d started at %v.", runID, measuredTimeUTC)

	run := models.Run{
		ID:                    runID,
		MeasuredTimeUTC:       measuredTimeUTC,
		ClusterCPUMillicore:   clusterCPUMillieCores,
		ClusterMemoryMegaByte: clusterMemoryMegaByte,
		Resources:             reqResources}
	err = sqlClient.SaveRequiredResources(run)
	if err != nil {
		return errors.WithMessage(err, "error saving resources")
	}

	log.Infof("Run %d finished successfully at %v", runID, time.Now().UTC())
	return nil
}
