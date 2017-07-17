package main

import (
	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

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

// CreateError creates a new error with fields filled.
func CreateError(etype ErrorType, code int, message string) Error {
	if etype != ErrorTypePrivate && etype != ErrorTypePublic {
		panic("Error type needs to be either ErrorTypePrivate or ErrorTypePublic")
	} else {
		return Error{Type: etype, Code: code, Message: message}
	}
}

// LogPrivateError sends private errors to sentry
// for internal error logging.
func LogPrivateError(etype ErrorType, e error) {
	if etype != ErrorTypePrivate {
		panic("Error type needs to be ErrorTypePrivate for private error logging")
	}
	raven.CaptureError(e, nil)
}

// LogPublicError returns error response to user
func LogPublicError(c *gin.Context, etype ErrorType, code int, message string) {
	if etype != ErrorTypePublic {
		panic("Error type needs to be ErrorTypePublic for public error logging")
	}
	c.JSON(code, gin.H{"error": gin.H{"code": code, "message": message}})
}
