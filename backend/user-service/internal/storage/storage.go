package storage

import (
	"errors"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")

	ErrRefreshNotValid = errors.New("refresh token expired or not found")

	ErrSessionNotFound = errors.New("session not found")
)
