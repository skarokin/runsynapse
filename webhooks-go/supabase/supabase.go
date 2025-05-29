package supabase

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

type SupabaseClient struct {
    Pool *pgxpool.Pool
}

func NewClient(databaseURL string) (*SupabaseClient, error) {
    if databaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL environment variable is required")
    }

    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse database URL: %w", err)
    }

    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }

    if err := pool.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return &SupabaseClient{Pool: pool}, nil
}

func (s *SupabaseClient) Close() {
    s.Pool.Close()
}