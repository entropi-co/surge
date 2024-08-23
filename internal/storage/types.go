package storage

import (
	"database/sql"
	"surge/internal/schema"
	"time"
)

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

func NullStringToPointer(str sql.NullString) *string {
	if str.Valid {
		return &str.String
	}
	return nil
}

func NullTimeToPointer(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func NullableToPointer[T any](container sql.Null[T]) *T {
	if container.Valid {
		return &container.V
	}
	return nil
}
