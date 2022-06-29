package teal

import (
	"errors"
)

var (
	ErrDoesNotExist = errors.New("the item does not exist")
	ErrNoRows       = errors.New("no items found")

	ErrDuplicateUsername = errors.New("username already exists")

	ErrNoAuthHeader  = errors.New("no authentication headers")
	ErrInvalidCreds  = errors.New("invalid credentials")
	ErrAPIKeyExpired = errors.New("api key expired")
)
