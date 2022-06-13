package teal

import (
	"errors"
	"fmt"
)

var ErrDoesNotExist = errors.New("the item does not exist")
var ErrNoRows = errors.New("no items found")

type ValidationError struct {
	Message string `json:"message"`
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{fmt.Sprintf("%s: %s", field, message)}
}
