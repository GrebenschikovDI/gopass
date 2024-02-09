package users

import (
	"GoPass/internal/server"
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo Repo
}

func NewUseCase(repo Repo) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) RegisterUser(ctx context.Context, username, password string) (*User, error) {
	existingUser, err := u.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, errors.Wrapf(server.ErrUserExists, "register user err, username: %s exists", username)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(err, "register user err, username: %s", username)
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, errors.Wrapf(err, "register user err, can't hash password")
	}

	newUser := &User{
		Login:    username,
		Password: passwordHash,
	}

	if err := u.repo.Create(ctx, newUser); err != nil {
		return nil, errors.Wrapf(err, "register user err, can't create user %s", username)
	}

	user, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrapf(err, "register user err, can't get user %s", username)
	}
	return user, nil
}

func (u *UseCase) AuthenticateUser(ctx context.Context, username, password string) (*User, error) {
	user, err := u.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrapf(err, "auth user err, can't get username: %s", username)
	}
	if user == nil {
		return nil, errors.Wrapf(server.ErrUnauthorized, "auth user err, username: %s", username)
	}
	if err := comparePasswords(user.Password, password); err != nil {
		return nil, errors.Wrapf(err, "auth user err, probably wrong password, username: %s", username)
	}
	return user, nil
}

func (u *UseCase) GetByUsername(ctx context.Context, username string) (*User, error) {
	existingUser, err := u.repo.GetByUsername(ctx, username)

	if err != nil {
		return nil, errors.Wrapf(err, "get user by username err, username %s", username)
	}
	if existingUser == nil {
		return nil, errors.Wrapf(server.ErrUserNotFound, "get user by username err, username %s", username)
	}
	return existingUser, nil
}

func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func comparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return server.ErrUnauthorized
	}
	return nil
}
