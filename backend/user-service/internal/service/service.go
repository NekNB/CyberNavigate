package service

import "errors"

var (
	AuthenticationError = errors.New("Authorization Error: Incorrect Login or Password")
)
