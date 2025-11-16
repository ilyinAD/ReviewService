package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewPgxPoolConfig(cfg *DBConfig) (*pgxpool.Config, error) {
	fmt.Println(cfg.ToDSN())

	connConfig, err := pgxpool.ParseConfig(cfg.ToDSN())
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewPgxPool: %w", err)
	}

	connConfig.MaxConns = 12
	connConfig.MinConns = 2
	connConfig.MaxConnIdleTime = time.Second

	return connConfig, nil
}

func NewPgxPool(connConfig *pgxpool.Config, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	lc.Append(fx.Hook{OnStop: func(_ context.Context) error {
		log.Println("Closing database pool")
		pool.Close()
		return nil
	}})

	return pool, nil
}
