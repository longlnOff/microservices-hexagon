package postgres

import (
	"context"
	"log"
	"embed"
	"fmt"

	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DB struct {
	*pgxpool.Pool
	url string
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

	return &DB{db, connectionString}, nil
}

func (db *DB) MigrateUpTo(version int) error {
	goose.SetBaseFS(embedMigrations)

	// Create a pgx connection config
	connConfig, err := pgx.ParseConfig(db.url)
	if err != nil {
		log.Fatalf("Failed to parse connection string: %v", err)
	}
	
	// Register the pgx driver with the database/sql package
	_ = stdlib.RegisterConnConfig(connConfig)
	
	// Open a database/sql connection using the registered driver
	database := stdlib.OpenDB(*connConfig)
	defer database.Close()

	// Set the database dialect for goose
	goose.SetDialect("postgres")

	// Get database status
	err = goose.Status(database, "migrations")
	if err != nil {
		log.Fatalf("Failed to get database status: %v", err)
	}
	
	// Run migrations up to the most recent version
	err = goose.UpTo(database, "migrations", int64(version))
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	
	log.Println("Migrations completed successfully")

	return nil
}

// TODO: reference: https://stackoverflow.com/questions/76865674/how-to-use-goose-migrations-with-pgx
// TODO: reference: https://github.com/bagashiz/go-pos/blob/main/internal/adapter/storage/postgres/db.go

func (db *DB) ErrorCode(err error) string {
	pgErr := err.(*pgconn.PgError)
	return pgErr.Code
}

