package repository

import (
	"context"
	"database/sql"

	"time"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/equinor/radix-cost-allocation/pkg/models"
	"github.com/pkg/errors"
)

const (
	defaultQueryTimeout  int    = 30
	nodeTvpTypeName      string = "cost.node_upsert_type"
	containerTvpTypeName string = "cost.container_upsert_type"
)

type sqlRepository struct {
	db           *sql.DB
	queryTimeout int
	ctx          context.Context
}

func NewSqlRepository(db *sql.DB, QueryTimeout int, context context.Context) Repository {
	return &sqlRepository{
		db:           db,
		queryTimeout: QueryTimeout,
		ctx:          context,
	}
}

// BulkUpsertContainers writes the list of containers to the database
// If the container exists in the database, only LastKnowRunningAt will be updated
func (repo *sqlRepository) BulkUpsertContainers(containers []models.ContainerBulkTvp) error {
	nodeArg := sql.Named("containers",
		mssql.TVP{
			TypeName: containerTvpTypeName,
			Value:    containers,
		},
	)
	if err := repo.executeWithTransaction("exec cost.container_upsert_bulk @containers", nodeArg); err != nil {
		return errors.WithMessage(err, "BulkUpsertContainers")
	}
	return nil
}

func (repo *sqlRepository) BulkUpsertNodes(nodes []models.NodeBulkTvp) error {
	nodeArg := sql.Named("nodes",
		mssql.TVP{
			TypeName: nodeTvpTypeName,
			Value:    nodes,
		},
	)
	if err := repo.executeWithTransaction("exec cost.node_upsert_bulk @nodes", nodeArg); err != nil {
		return errors.WithMessage(err, "BulkUpsertNodes")
	}
	return nil
}

func (repo *sqlRepository) createConnection() (*sql.Conn, error) {
	ctx, cancel := repo.getContext()
	defer cancel()
	return repo.db.Conn(ctx)
}

func (repo *sqlRepository) executeWithTransaction(query string, args ...interface{}) error {
	ctx, cancel := repo.getContext()
	defer cancel()

	conn, err := repo.createConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repo *sqlRepository) getContext() (ctx context.Context, cancel context.CancelFunc) {
	if repo.queryTimeout > 0 {
		return context.WithTimeout(repo.ctx, time.Duration(repo.queryTimeout)*time.Second)
	} else {
		return context.WithCancel(repo.ctx)
	}
}
