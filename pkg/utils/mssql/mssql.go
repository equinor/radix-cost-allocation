package utils

import (
	"database/sql"
	"fmt"

	mssql "github.com/denisenkom/go-mssqldb"
)

func OpenSqlServer(server, database, userID, password string, port int) (*sql.DB, error) {
	c, err := mssql.NewConnector(
		GetSqlServerDsn(server, database, userID, password, port),
	)
	if err != nil {
		return nil, err
	}

	return sql.OpenDB(c), nil
}

func GetSqlServerDsn(server, database, userID, password string, port int) string {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
		server, userID, password, database)

	if port > 0 {
		dsn = fmt.Sprintf("%s;port=%d", dsn, port)
	}

	return dsn
}
