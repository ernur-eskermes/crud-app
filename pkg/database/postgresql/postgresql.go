package postgresql

import (
	"context"
	"log"
	"time"

	"github.com/ernur-eskermes/crud-app/pkg/repeatable"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type StorageConfig struct {
	ConnStr string
	Logger  pgx.Logger
}

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, sc StorageConfig) (pool *pgxpool.Pool, err error) {
	config, err := pgxpool.ParseConfig(sc.ConnStr)
	if err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}

	config.ConnConfig.Logger = sc.Logger

	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.ConnectConfig(ctx, config)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return pool, nil
}
