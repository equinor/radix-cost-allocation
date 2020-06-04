package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

// todo! create write only connection string? dont need read/admin access
const port = 1433

func main() {
	sqlClient := Init(os.Getenv("SQL_SERVER"), os.Getenv("SQL_DATABASE"), port, os.Getenv("SQL_USER"), os.Getenv("SQL_PASSWORD"))
	promClient := PrometheusClient{Address: os.Getenv("PROMETHEUS_API")}

	measuredTimeUTC := time.Now().UTC()
	reqResources, err := promClient.GetRequiredResources(measuredTimeUTC)
	if err != nil {
		log.Fatal("Error getting required resources: ", err.Error())
	}

	runID, err := sqlClient.SaveRun(measuredTimeUTC)
	if err != nil {
		log.Fatal("Error creating Run: ", err.Error())
	}
	fmt.Printf("Run %d started at %v.\n", runID, measuredTimeUTC)

	run := Run{ID: runID, MeasuredTimeUTC: measuredTimeUTC, Resources: reqResources}
	err = sqlClient.SaveRequiredResources(run)
	if err != nil {
		log.Fatal("Error saving resources: ", err.Error())
	}

	fmt.Printf("Run %d finished successfully at %v", runID, time.Now().UTC())
	sqlClient.Close()
}
