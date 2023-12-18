package model

import "errors"

type User struct {
	ID
	Name     string
	Email    string
	Password string
}

var (
	ErrUserNotExist = errors.New("User with this email and password not found")
	ErrUserExist    = errors.New("User with this email already exists")
)
