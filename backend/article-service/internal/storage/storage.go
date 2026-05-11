package storage

import (
	"errors"
)

var (
	ErrArticleExists            = errors.New("article already exists")
	ErrArticleNotFound          = errors.New("article not found")
	ErrArticleTextNotCreatedYet = errors.New("article text not created yet")
)
