package models

import "time"

type TokensDTO struct {
	AccessToken  string
	RefreshToken string
}

type UserDTO struct {
	UserId    string
	Username  string
	IsAdmin   bool
	CreatedAt time.Time
}

type SessionDTO struct {
	SessionId    string
	RefreshToken string
	User         UserDTO
}
