package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

// NewSQLRepository returns a SQL implementation of the Repository interface
func NewSQLRepository(context context.Context, db *sql.DB, QueryTimeout int) Repository {
	return &sqlRepository{
		db:           db,
		queryTimeout: QueryTimeout,
		ctx:          context,
	}
}

// BulkUpsertContainers writes the list of containers to the database
func (repo *sqlRepository) BulkUpsertContainers(containers []ContainerBulkDto) error {
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

// BulkUpsertNodes writes the list of nodes to the database
func (repo *sqlRepository) BulkUpsertNodes(nodes []NodeBulkDto) error {
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
	defer func() {
		if err = conn.Close(); err != nil {
			log.Error().Err(err).Msg("Unable to close connection")
		}
	}()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w (transaction error: %w)", rollbackErr, err)
		}
		return err
	}

	return tx.Commit()
}

func (repo *sqlRepository) getContext() (ctx context.Context, cancel context.CancelFunc) {
	if repo.queryTimeout > 0 {
		return context.WithTimeout(repo.ctx, time.Duration(repo.queryTimeout)*time.Second)
	}

	return context.WithCancel(repo.ctx)
}
