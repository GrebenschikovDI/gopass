package users

import "context"

type Repo interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, userID int) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
