package config

type AppConfig struct {
	PrometheusAPI    string
	CronSchedule     string `envconfig:"default=0 0 * * * *"`
	PodSyncSchedule  string `envconfig:"default=10 * * * * *"`
	NodeSyncSchedule string `envconfig:"default=0 0/5 * * * *"`
	SQL              SQLConfig
}

type SQLConfig struct {
	Server       string
	Database     string `envconfig:"default=sqldb-radix-cost-allocation"`
	User         string
	Password     string
	Port         int `envconfig:"default=1433"`
	QueryTimeout int `envconfig:"default=30"`
}
