package main

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/v3/models"
	"log"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

// todo! create write only connection string? dont need read/admin access
const port = 1433

func main() {
	promClient := PrometheusClient{Address: os.Getenv("PROMETHEUS_API")}
	sqlClient := NewSQLClient(os.Getenv("SQL_SERVER"), os.Getenv("SQL_DATABASE"), port, os.Getenv("SQL_USER"), os.Getenv("SQL_PASSWORD"))

	moveResourceRequestsFromPrometheusToSQLDB(promClient, sqlClient)
	// printCostBetweenDates(time.Now().UTC().AddDate(0, 0, -3), time.Now().UTC(), promClient, sqlClient)

	sqlClient.Close()
}

func printCostBetweenDates(from, to time.Time, promClient PrometheusClient, sqlClient SQLClient) {
	runs, err := sqlClient.GetRunsBetweenTimes(from, to)
	if err != nil {
		log.Fatal("Error getting runs: ", err.Error())
	}

	cost := models.NewCost(from, to, runs)
	costJSON, err := json.Marshal(cost.Applications)
	if err != nil {
		log.Fatal("Error converting to json: ", err.Error())
	}
	fmt.Println(string(costJSON))
}

func moveResourceRequestsFromPrometheusToSQLDB(promClient PrometheusClient, sqlClient SQLClient) {
	measuredTimeUTC := time.Now().UTC()
	reqResources, err := promClient.GetRequiredResources(measuredTimeUTC)
	if err != nil {
		log.Fatal("Error getting required resources: ", err.Error())
	}

	clusterCPUCores, err := promClient.GetClusterTotalCPUCoresFromPrometheus(measuredTimeUTC)
	if err != nil {
		log.Fatal("Error getting node cpu count: ", err.Error())
	}
	clusterCPUMillieCores := clusterCPUCores * 1000

	clusterMemoryBytes, err := promClient.GetClusterTotalMemoryBytesFromPrometheus(measuredTimeUTC)
	if err != nil {
		log.Fatal("Error getting node cpu count: ", err.Error())
	}
	clusterMemoryMegaByte := clusterMemoryBytes / 1000000

	runID, err := sqlClient.SaveRun(measuredTimeUTC, clusterCPUMillieCores, clusterMemoryMegaByte)
	if err != nil {
		log.Fatal("Error creating Run: ", err.Error())
	}
	fmt.Printf("Run %d started at %v.\n", runID, measuredTimeUTC)

	run := models.Run{
		ID:                    runID,
		MeasuredTimeUTC:       measuredTimeUTC,
		ClusterCPUMillicore:   clusterCPUMillieCores,
		ClusterMemoryMegaByte: clusterMemoryMegaByte,
		Resources:             reqResources}
	err = sqlClient.SaveRequiredResources(run)
	if err != nil {
		log.Fatal("Error saving resources: ", err.Error())
	}

	fmt.Printf("Run %d finished successfully at %v", runID, time.Now().UTC())
}
