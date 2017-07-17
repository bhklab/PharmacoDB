package main

import raven "github.com/getsentry/raven-go"

// ErrorType is one of
// private or public (0 or 1).
type ErrorType int

// Error types.
const (
	ErrorTypePrivate ErrorType = 0
	ErrorTypePublic  ErrorType = 1
)

// Error is an internal error structure.
type Error struct {
	Error   error
	Type    ErrorType
	Code    int
	Message string
}

// Sentry DSN for internal error logging.
func init() {
	raven.SetDSN("https://71d8d1bc8e4843eeba979fdaadebe48b:df30d2048fc44b5185809f04ba9d2294@sentry.io/186627")
}

// NewError returns an empty error
func NewError() Error {
	return Error{}
}
