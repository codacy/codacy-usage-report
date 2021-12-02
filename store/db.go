package store

import (
	"database/sql"
	"fmt"

	"github.com/codacy/codacy-usage-report/config"
	_ "github.com/lib/pq"
)

func dbConnectionString(dbConfig config.DatabaseConfiguration) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s port=%d host=%s sslmode=%s", dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.Host, dbConfig.SslMode)
}

func connectToDB(dbConfig config.DatabaseConfiguration) (*sql.DB, error) {
	return sql.Open("postgres", dbConnectionString(dbConfig))
}
