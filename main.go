package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/equinor/radix-cost-allocation/config"
	"github.com/equinor/radix-cost-allocation/run"
	log "github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
)

var Config config.AppConfig

func main() {
	if err := envconfig.Init(&Config); err != nil {
		log.Fatal(err)
	}

	stopCh := make(chan struct{})
	go run.InitAndRunOldDataCollector(Config.PrometheusAPI, Config.CronSchedule, Config.SQL, stopCh)
	go func() {
		if err := run.InitAndStartCollector(Config, stopCh); err != nil {
			log.Fatal(err)
		}
	}()

	sigTerm := make(chan os.Signal, 1)
	go func() {
		signal.Notify(sigTerm, syscall.SIGTERM)
		signal.Notify(sigTerm, syscall.SIGINT)
	}()

	<-sigTerm
	close(stopCh)
}
