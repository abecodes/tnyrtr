package tnyrtr

import (
	"errors"
)

// shutdownError is a type used to help with graceful app shutdown
type shutdownError struct {
	Message string
}

// NewShutdownErr returns an error that will cause an graceful app shutdown
func NewShutdownErr(msg string) error {
	return &shutdownError{Message: msg}
}

// Error implements the error interface
func (se *shutdownError) Error() string {
	return se.Message
}

// IsShutdown checks if an error is a shutdownErr
func IsShutdown(err error) bool {
	var se *shutdownError
	return errors.As(err, &se)
}
