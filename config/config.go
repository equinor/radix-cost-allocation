package config

// AppConfig holds all configuration options for the application
type AppConfig struct {
	PrometheusAPI      string
	CronSchedule       string `envconfig:"default=0 0 * * * *"`
	Schedule           CronSchedule
	SQL                SQLConfig
	AppNameExcludeList []string `envconfig:"optional"`
	LogLevel           string   `envconfig:"default=info"`
}

// SQLConfig defines configuration settings used to manage connections to SQL Server
type SQLConfig struct {
	Server       string
	Database     string `envconfig:"default=sqldb-radix-cost-allocation"`
	User         string
	Password     string
	Port         int `envconfig:"default=1433"`
	QueryTimeout int `envconfig:"default=30"`
}

// CronSchedule defines cron schedules for jobs
type CronSchedule struct {
	PodSync  string `envconfig:"default=10 * * * * *"`
	NodeSync string `envconfig:"default=0 0/5 * * * *"`
}
