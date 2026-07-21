package services

import "errors"

var (
	SessionAlreadyExists error = errors.New("Session for current user already exists")
)
