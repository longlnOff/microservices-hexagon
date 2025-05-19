package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"log"

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
	DB *pgxpool.Pool
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
	config.HealthCheckPeriod = 2 * time.Second

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}


	dbObject := DB{
		db,
		connectionString,
	}
	err = dbObject.checkConnection(context.Background())
	if err != nil {
		return nil, err
	}
	go dbObject.startHealthCheck(context.Background())

	return &dbObject, nil
}

func (db *DB) startHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(db.DB.Config().HealthCheckPeriod)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if err := db.checkConnection(ctx); err != nil {
				log.Printf("DB health check failed: %v", err)
				// Force reconnection attempt
				db.reconnect(ctx)
			}
		case <-ctx.Done():
			log.Println("Stopping DB health check")
			return
		}
	}
}

// checkConnection verifies if the database connection is still alive
func (dm *DB) checkConnection(ctx context.Context) error {
	// Create timeout context for the ping
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	
	// Use Ping to verify connection
	err := dm.DB.Ping(pingCtx)
	if err != nil {
		return err
	}
	
	// Optionally perform a simple query to further verify connection
	conn, err := dm.DB.Acquire(pingCtx)
	if err != nil {
		return err
	}
	defer conn.Release()
	
	var result int
	err = conn.QueryRow(pingCtx, "SELECT 1").Scan(&result)
	if err != nil {
		return err
	}
	
	return nil
}

// reconnect attempts to reconnect to the database by closing and recreating the pool
func (dm *DB) reconnect(ctx context.Context) {
	log.Println("Attempting to reconnect to database...")
	
	// Close existing pool if it exists
	if dm != nil {
		dm.DB.Close()
	}
	db, err := pgxpool.NewWithConfig(ctx, dm.DB.Config())
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err = db.Ping(ctx)
	if err != nil {
		log.Printf("Failed to reconnect: %v", err)
	} else {
		dm.DB = db
		log.Println("Successfully reconnected to database")
	}
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

// ErrorCode extracts the Postgres error code from different types of pgx errors
func (db *DB) ErrorCode(err error) string {
    if err == nil {
        return ""
    }
    
    // Handle PgError (for query errors, constraint violations, etc.)
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        return pgErr.Code
    }
    
    // Handle ConnectError (for connection failures)
    var connectErr *pgconn.ConnectError
    if errors.As(err, &connectErr) {
        // ConnectError doesn't have a Code field directly
        // So we'll return a custom code for connection errors
        return "connection_error"
    }
    
    // Try to check for errors implementing the PostgreSQL error interface
    type pgxErrorCode interface {
        SQLState() string
    }
    
    var sqlStateErr pgxErrorCode
    if errors.As(err, &sqlStateErr) {
        return sqlStateErr.SQLState()
    }
    
    // Additional checks for network-related errors
    if errors.Is(err, context.DeadlineExceeded) {
        return "timeout_error"
    }
    
    if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
        return "connection_closed"
    }
    
    // Default case - unknown error type
    return "unknown_error"
}
