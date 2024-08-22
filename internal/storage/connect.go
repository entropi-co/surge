package storage

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"surge/internal/conf"
	"surge/internal/schema"

	_ "github.com/jackc/pgx/v5"
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

func CreateQueries(conn *sql.DB) *schema.Queries {
	return schema.New(conn)
}

func NewString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: true}
}

func NewStringNull() sql.NullString {
	return sql.NullString{Valid: false}
}

func NewNullableString(p *string) sql.NullString {
	if p == nil {
		return NewStringNull()
	} else {
		return NewString(*p)
	}
}
