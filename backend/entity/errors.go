package entity

import "errors"

var (
	ErrNotFound           = errors.New("data not found")
	ErrAlreadyExists      = errors.New("data already exists")
	ErrInvalidID          = errors.New("invalid id")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoChanges          = errors.New("no changes detected")
)
