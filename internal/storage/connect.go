package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"surge/internal/conf"
)

func CreateDatabaseConnection(config *conf.SurgeDatabaseConfigurations) *sql.DB {
	c, err := sql.Open("postgres", config.Url)
	if err != nil {
		logrus.WithError(err).Fatalf("Unable to connect to the database\n")
	}

	return c
}

func CloseDatabase(conn *sql.DB) {
	err := conn.Close()
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to close the database connection\n")
	}
}
