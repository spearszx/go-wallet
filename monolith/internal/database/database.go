package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spearszx/go-wallet/monolith/internal/config"
)

var (
	ErrMaxRetries = errors.New("max retries exceeded, giving up")
)

var (
	defaultTimeout  = time.Second * 5
	defaultAttempts = 3
)

func NewWithRetries(config *config.DbOptions) (*pgxpool.Pool, error) {
	ticker := time.NewTicker(defaultTimeout)
	defer ticker.Stop()

	for attempt := 0; attempt < config.DbMaxRetries; attempt++ {
		fmt.Printf("attempting to connect to database №%d\n", attempt+1)

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		pool, err := pgxpool.New(ctx, config.ConnectionString())
		if err == nil {
			if err = pool.Ping(ctx); err == nil {
				return pool, nil
			}
		}

		fmt.Println(err)

		<-ticker.C
	}

	return nil, fmt.Errorf("end")
}
