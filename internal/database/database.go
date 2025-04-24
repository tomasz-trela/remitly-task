package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("could not open db connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping db: %v", err)
	}

	DB = db
	return DB, nil
}
