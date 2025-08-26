package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"my-habr/services/auth/db"
	"my-habr/services/auth/model"
)

type UserRepository struct {
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := db.Pool.Exec(ctx,
		"INSERT INTO users(email, password) VALUES ($1, $2)",
		user.Email, user.Password)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	log.Printf("🔍 Querying DB for email: %s", email)
	row := db.Pool.QueryRow(ctx,
		"SELECT id, email, password FROM users WHERE email = $1", email)

	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("🔍 User not found in DB: %s", email)
			return nil, nil
		}
		log.Printf("❌ DB query error: %v", err)
		return nil, err
	}

	log.Printf("✅ User found in DB: ID=%d, Email=%s", user.ID, user.Email)
	return &user, nil
}
