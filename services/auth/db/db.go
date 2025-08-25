// db/db.go
package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var Pool *pgxpool.Pool

func InitDB(dsn string) error {
	var err error
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}

	// Повторяем попытки подключения
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		Pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			// Проверяем подключение
			if err = Pool.Ping(context.Background()); err == nil {
				fmt.Println("✅ PostgreSQL connected")
				return nil
			}
		}

		fmt.Printf("❌ PostgreSQL not ready, retrying %d/%d: %v\n", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to connect to PostgreSQL after %d retries", maxRetries)
}
