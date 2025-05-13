package postgres

import (
	"context"
	"fmt"

	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
	// "github.com/jackc/pgx/v5/stdlib"
	// "github.com/pressly/goose/v3"
)

type DB struct {
	*pgxpool.Pool
}

func New(
	connectionString string,
	maxOpenConns int,
	minIdleConns int,
	maxIdleTime time.Duration,
) (
	*DB,
	error,
) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		return nil, nil
	}

	// Set the maximum number of connections in the pool
	config.MaxConns = int32(maxOpenConns)
	config.MaxConnIdleTime = maxIdleTime
	config.MinIdleConns = int32(minIdleConns)

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// TODO: add migrations with goose in code & set configuration in environment to specify migration to version
// TODO: reference: https://stackoverflow.com/questions/76865674/how-to-use-goose-migrations-with-pgx
// TODO: reference: https://github.com/bagashiz/go-pos/blob/main/internal/adapter/storage/postgres/db.go

func (db *DB) ErrorCode(err error) string {
	pgErr := err.(*pgconn.PgError)
	fmt.Println(pgErr.Message)
	return pgErr.Code
}
