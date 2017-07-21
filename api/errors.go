package api

import (
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

// Sentry DSN for internal error logging.
func init() {
	raven.SetDSN("https://71d8d1bc8e4843eeba979fdaadebe48b:df30d2048fc44b5185809f04ba9d2294@sentry.io/186627")
}

// LogPrivateError sends private errors to sentry for internal error logging.
func LogPrivateError(err error) {
	raven.CaptureError(err, nil)
}

// LogPublicError responds with error message upon routine failure.
func LogPublicError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"code": code, "message": message})
}

// LogInternalServerError writes error message with status code http.StatusInternalServerError
func LogInternalServerError(c *gin.Context, message string) {
	LogPublicError(c, http.StatusInternalServerError, message)
}

// LogNotFoundError writes error message with status code http.StatusNotFound
func LogNotFoundError(c *gin.Context, message string) {
	LogPublicError(c, http.StatusNotFound, message)
}

// LogBadRequestError writes error message with status code http.StatusBadRequest
func LogBadRequestError(c *gin.Context, message string) {
	LogPublicError(c, http.StatusBadRequest, message)
}
