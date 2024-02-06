package server

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrUnauthorized  = errors.New("authentication failed")
	ErrAlreadyExists = errors.New("order already exists")
	ErrAlreadyTaken  = errors.New("order is taken by another user")
	ErrEmptyField    = errors.New("username or password must not be empty")
	ErrBadFormat     = errors.New("bad format of order")
	ErrNotValid      = errors.New("order not valid")
)
