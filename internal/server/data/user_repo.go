package data

import (
	"GoPass/internal/server/users"
	"context"
	"database/sql"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) users.Repo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, user *users.User) error {
	_, err := u.db.ExecContext(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2)",
		user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepo) GetByID(ctx context.Context, id int) (*users.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE id = $1", id)
	user := &users.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepo) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE username = $1",
		username)
	user := &users.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
