package inits

import (
    "context"
    "fmt"
	"time"

	"github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

func NewSupabaseClient(databaseURL string) (*pgxpool.Pool, error) {
    if databaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL environment variable is required")
    }

    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol // disable prepared statement caching to avoid conflicts
	config.MaxConns = 3 
	config.MinConns = 0 
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 1 * time.Minute

    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }

    if err := pool.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return pool, nil
}