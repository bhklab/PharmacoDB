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

// Sentry DSN for internal error logging.
func init() {
	raven.SetDSN("https://71d8d1bc8e4843eeba979fdaadebe48b:df30d2048fc44b5185809f04ba9d2294@sentry.io/186627")
}

// LogPrivateError sends private errors to sentry
// for internal error logging.
func LogPrivateError(typ ErrorType, err error) {
	if typ != ErrorTypePrivate {
		panic("Error type needs to be ErrorTypePrivate for private error logging")
	}
	raven.CaptureError(err, nil)
}

// LogPublicError returns error response to user
func LogPublicError(c *gin.Context, typ ErrorType, code int, message string) {
	if typ != ErrorTypePublic {
		panic("Error type needs to be ErrorTypePublic for public error logging")
	}
	c.JSON(code, gin.H{"code": code, "message": message})
}
