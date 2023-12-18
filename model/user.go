package model

import "errors"

type User struct {
	ID       uint64
	Name     string
	Email    string
	Password string
}

var (
	ErrLoginMismatch       = errors.New("User with this email and password not found")
	ErrUserWithEmailExists = errors.New("User with this email already exists")
)
