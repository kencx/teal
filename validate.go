package teal

import "fmt"

type ValidationError struct {
	Message string `json:"message"`
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{fmt.Sprintf("%s: %s", field, message)}
}
