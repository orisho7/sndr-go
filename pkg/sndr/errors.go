package sndr

import (
	"fmt"
)

// APIError represents a structured error returned by the SNDR API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Fields     map[string]string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("sndr api error (status %d): %s: %s", e.StatusCode, e.Code, e.Message)
}

// Is reports whether err is an APIError.
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}
