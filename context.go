package tnyrtr

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// ctxKey represents the valuetype for context key
// custom key type so it can't be overwritten
type ctxKey int

// key is how request values are stored/retrieved
const key ctxKey = 1

// Values represent state for each request
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

var errMissing = errors.New("values missing in context")

// GetValues returns the state values from the context
func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(key).(*Values)

	if !ok {
		return nil, errMissing
	}

	return v, nil
}

// GetTraceID returns the traceID from the context
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*Values)

	if !ok {
		return "0000-0000-0000-0000"
	}

	return v.TraceID
}

// GetTime returns the UTC time from the context
func GetTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(key).(*Values)

	if !ok {
		return time.Now().UTC()
	}

	return v.Now
}

// GetStatus returns the status code for the request from the context
func GetStatus(ctx context.Context) int {
	v, ok := ctx.Value(key).(*Values)

	if !ok {
		return http.StatusInternalServerError
	}

	return v.StatusCode
}

// SetStatusCode sets the satus code for the context
func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(key).(*Values)

	if !ok {
		return errMissing
	}

	v.StatusCode = statusCode

	return nil
}
