// db/db.go
package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var Pool *pgxpool.Pool

func InitBD(dsn string) error {
	var err error
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil
	}

	config.MaxConns = 20
	config.HealthCheckPeriod = 30 * time.Second

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	if err = Pool.Ping(context.Background()); err != nil {
		return err
	}

	fmt.Println("PostgreSQL connected")
	return nil
}
