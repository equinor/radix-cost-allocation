package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/equinor/radix-cost-allocation/models"
)

// SQLClient used to perform sql queries
type SQLClient struct {
	Server   string
	Port     int
	Database string
	UserID   string
	Password string
	db       *sql.DB
}

// NewSQLClient create sqlclient and setup the db connection
func NewSQLClient(server, database string, port int, userID, password string) SQLClient {
	sqlClient := SQLClient{
		Server:   server,
		Database: database,
		Port:     port,
		UserID:   userID,
		Password: password,
	}
	sqlClient.db = sqlClient.setupDBConnection()
	return sqlClient
}

// SaveRequiredResources inserts all required resources under run.Resources
func (sqlClient SQLClient) SaveRequiredResources(run models.Run) error {
	tsql := `INSERT INTO cost.required_resources (run_id, wbs, application, environment, component, cpu_millicores, memory_mega_bytes, replicas) 
	VALUES (@runId, @wbs, @application, @environment, @component, @cpuMillicores, @memoryMegaBytes, @replicas); select convert(bigint, SCOPE_IDENTITY());`
	for _, req := range run.Resources {
		_, err := sqlClient.execSQL(tsql,
			sql.Named("runId", run.ID),
			sql.Named("wbs", req.WBS),
			sql.Named("application", req.Application),
			sql.Named("environment", req.Environment),
			sql.Named("component", req.Component),
			sql.Named("cpuMillicores", req.CPUMillicore),
			sql.Named("memoryMegaBytes", req.MemoryMegaBytes),
			sql.Named("replicas", req.Replicas))
		if err != nil {
			return fmt.Errorf("Failed to insert req resources %v", err)
		}
	}
	return nil
}

// SaveRun inserts a new run, returns id
func (sqlClient SQLClient) SaveRun(measuredTime time.Time, clusterCPUMillicores int) (int64, error) {
	tsql := "INSERT INTO cost.runs (measured_time_utc, cluster_cpu_millicores) VALUES (@measuredTimeUTC, @clusterCPUMillicores); select convert(bigint, SCOPE_IDENTITY());"
	return sqlClient.execSQL(tsql,
		sql.Named("measuredTimeUTC", measuredTime),
		sql.Named("clusterCPUMillicores", clusterCPUMillicores))
}

// Close the underlying db connection
func (sqlClient SQLClient) Close() {
	sqlClient.db.Close()
}

// SetupDBConnection sets up db connection
func (sqlClient SQLClient) setupDBConnection() *sql.DB {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		sqlClient.Server, sqlClient.UserID, sqlClient.Password, sqlClient.Port, sqlClient.Database)

	var err error

	// Create connection pool
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
	return db
}

func (sqlClient SQLClient) execSQL(tsql string, args ...interface{}) (int64, error) {
	ctx, err := sqlClient.verifyConnection()
	if err != nil {
		return -1, err
	}

	stmt, err := sqlClient.db.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}

func (sqlClient SQLClient) verifyConnection() (context.Context, error) {
	ctx := context.Background()
	var err error

	if sqlClient.db == nil {
		err = errors.New("CreateRun: db is null")
		return ctx, err
	}

	// Check if database is alive.
	err = sqlClient.db.PingContext(ctx)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}
