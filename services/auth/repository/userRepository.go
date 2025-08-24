package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"my-habr/services/auth/db"
	"my-habr/services/auth/model"
)

func main() {

}

type UserRepository struct {
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := db.Pool.Exec(ctx,
		"INSERT INTO users(email, password) VALUES ($1, $2)",
		user.EMail, user.Password)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row := db.Pool.QueryRow(ctx, "SELECT id, email, password FROM users WHERE email =$1", email)

	var user model.User
	err := row.Scan(&user.ID, &user.EMail, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
