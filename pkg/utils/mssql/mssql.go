package utils

import (
	"database/sql"
	"fmt"

	mssql "github.com/microsoft/go-mssqldb"
)

// OpenSQLServer opens a new connection to a SQL Server
func OpenSQLServer(server, database, userID, password string, port int) (*sql.DB, error) {
	c, err := mssql.NewConnector(
		GetSQLServerDsn(server, database, userID, password, port),
	)
	if err != nil {
		return nil, err
	}

	return sql.OpenDB(c), nil
}

// GetSQLServerDsn builds a SQL Server specific DSN
func GetSQLServerDsn(server, database, userID, password string, port int) string {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
		server, userID, password, database)

	if port > 0 {
		dsn = fmt.Sprintf("%s;port=%d", dsn, port)
	}

	return dsn
}
