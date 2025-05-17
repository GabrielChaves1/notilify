package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	db *sqlx.DB
}

type Options struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func DefaultPostgresOptions() Options {
	return Options{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

func NewPostgresClient(connectionString string, options Options) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(options.MaxOpenConns)
	db.SetMaxIdleConns(options.MaxIdleConns)
	db.SetConnMaxLifetime(options.ConnMaxLifetime)
	db.SetConnMaxIdleTime(options.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("couldn't connect to postgres database: %w", err)
	}

	return db, nil
}
