package main

import (
	"time"

	"github.com/equinor/radix-cost-allocation/clients"
	"github.com/equinor/radix-cost-allocation/models"
	"github.com/pkg/errors"
	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
	"golang.org/x/sync/semaphore"

	_ "github.com/denisenkom/go-mssqldb"
)

var Conf struct {
	PrometheusAPI string
	CronSchedule  string `envconfig:"default=0 * * * *"`
	SQL           struct {
		Server   string
		Database string `envconfig:"default=sqldb-radix-cost-allocation"`
		User     string
		Password string
		Port     int `envconfig:"default=1433"`
	}
}

var sem = semaphore.NewWeighted(1)

func main() {
	if err := envconfig.Init(&Conf); err != nil {
		log.Fatal(err)
	}

	promClient := clients.PrometheusClient{Address: Conf.PrometheusAPI}
	sqlClient, err := clients.NewSQLClient(Conf.SQL.Server, Conf.SQL.Database, Conf.SQL.Port, Conf.SQL.User, Conf.SQL.Password)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sqlClient.Close()

	initAndRunDataCollector(promClient, sqlClient)
}

func initAndRunDataCollector(promClient clients.PrometheusClient, sqlClient clients.SQLClient) {
	log.Infof("Registering cron job using schedule %s", Conf.CronSchedule)
	c := cron.New()
	c.AddFunc(Conf.CronSchedule, func() {
		if !sem.TryAcquire(1) {
			return
		}
		defer sem.Release(1)

		if err := moveResourceRequestsFromPrometheusToSQLDB(promClient, sqlClient); err != nil {
			log.Error(err)
		}
	})

	c.Run()
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
