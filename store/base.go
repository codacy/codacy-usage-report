package store

import (
	"database/sql"
	"fmt"

	"github.com/codacy/codacy-usage-report/config"
)

type baseStore struct {
	db *sql.DB
}

func (store *baseStore) Connect(dbConfiguration config.DatabaseConfiguration) error {
	dbConnection, err := connectToDB(dbConfiguration)
	if err != nil {
		return fmt.Errorf("Could not connect to database: %s", err.Error())
	}

	if err = dbConnection.Ping(); err != nil {
		return fmt.Errorf("Could not connect to database: %s", err.Error())
	}

	store.db = dbConnection
	return nil
}

func (store *baseStore) Close() error {
	fmt.Println("Closing connection")
	return store.db.Close()
}
