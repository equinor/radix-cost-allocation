package utils

import (
	"database/sql"
	"fmt"

	"github.com/microsoft/go-mssqldb/azuread"
)

// OpenSQLServer opens a new connection to a SQL Server
func OpenSQLServer(server, database string, port int) (*sql.DB, error) {
	c, err := azuread.NewConnector(GetSQLServerDsn(server, database, port))
	if err != nil {
		return nil, err
	}

	return sql.OpenDB(c), nil
}

// GetSQLServerDsn builds a SQL Server specific DSN
func GetSQLServerDsn(server, database string, port int) string {
	dsn := fmt.Sprintf("server=%s;database=%s;fedauth=ActiveDirectoryDefault", server, database)

	if port > 0 {
		dsn = fmt.Sprintf("%s;port=%d", dsn, port)
	}

	return dsn
}
