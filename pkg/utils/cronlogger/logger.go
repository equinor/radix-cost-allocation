package cronlogger

import (
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type CronLogger struct {
	logger *zerolog.Logger
}

var _ cron.Logger = CronLogger{}

func New(logger *zerolog.Logger) CronLogger {
	return CronLogger{logger: logger}
}

func (c CronLogger) Info(msg string, v ...interface{}) {
	c.logger.Info().Msgf(msg, v...)
}

func (c CronLogger) Error(err error, msg string, v ...interface{}) {
	c.logger.Error().Err(err).Msgf(msg, v...)
}
