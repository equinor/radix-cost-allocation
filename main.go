package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/equinor/radix-cost-allocation/config"
	"github.com/equinor/radix-cost-allocation/run"
	_ "github.com/microsoft/go-mssqldb"
	log "github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
)

var appConfig config.AppConfig

func main() {
	if err := envconfig.Init(&appConfig); err != nil {
		log.Fatal(err)
	}

	logLevel, err := log.ParseLevel(appConfig.LogLevel)
	if err != nil {
		log.Warnf("Log level '%s' is not valid. Using 'info' level", appConfig.LogLevel)
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	stopCh := make(chan struct{})

	go func() {
		if err := run.InitAndStartCollector(appConfig.SQL, appConfig.Schedule, appConfig.AppNameExcludeList, stopCh); err != nil {
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
