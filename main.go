package main

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/equinor/radix-cost-allocation/config"
	"github.com/equinor/radix-cost-allocation/run"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"
)

var appConfig config.AppConfig

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGTERM)
	defer cancel()

	if err := envconfig.Init(&appConfig); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize config")
	}

	ctx, err := setupLogger(ctx, appConfig.LogLevel, appConfig.PrettyPrint)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize logger")
	}

	info, _ := debug.ReadBuildInfo()
	log.Info().Str("version", info.Main.Version).Msg("Starting")

	go func() {
		if err := run.InitAndStartCollector(ctx, appConfig.SQL, appConfig.Schedule, appConfig.AppNameExcludeList); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()

	<-ctx.Done()
}
func setupLogger(ctx context.Context, logLevel string, prettyPrint bool) (context.Context, error) {
	zerolog.DurationFieldUnit = time.Millisecond
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	zerolog.SetGlobalLevel(level)
	if prettyPrint {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly})
	}
	zerolog.DefaultContextLogger = &log.Logger
	ctx = log.Logger.WithContext(ctx)
	return ctx, nil
}
