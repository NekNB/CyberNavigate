package storage

import (
	"errors"
)

var (
	ErrAlreadyExists = errors.New("Entity already exists")

	ErrNotFound = errors.New("Entity not found")

	ErrDataConfllict = errors.New("Data Conflict")
)
