package repository

import (
	"time"
)

// ContainerBulkTvp defines an object with fields for bulk insert/update used by the Repository interface
// Fields in this struct is one-to-one mapped to the SQL Server cost.container_upsert_type table value type (TVP).
// Order of fields in the struct must match the order of fields in the SQL Server type
type ContainerBulkTvp struct {
	ContainerID          string
	ContainerName        string
	PodName              string
	ApplicationName      string
	EnvironmentName      string
	ComponentName        string
	Wbs                  string
	StartedAt            time.Time
	LastKnowRunningAt    time.Time
	CPURequestMillicores int64
	MemoryRequestBytes   int64
	NodeName             string
}
