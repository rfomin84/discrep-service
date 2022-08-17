package mysql_client

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func NewMysqlClient(host, port, username, password, database string) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		return nil, err
	}
	return db, nil
}
