package api

import (
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

// Error is a custom public error implementation.
type Error struct {
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	Description interface{} `json:"description,omitempty"`
}

// Sentry DSN for internal error logging.
func init() {
	raven.SetDSN("https://71d8d1bc8e4843eeba979fdaadebe48b:df30d2048fc44b5185809f04ba9d2294@sentry.io/186627")
}

// LogSentry submits private errors to sentry.
func LogSentry(err error) {
	raven.CaptureError(err, nil)
}

// BadRequest responds with error status code 400, Bad Request.
func BadRequest(c *gin.Context, description interface{}) {
	err := Error{http.StatusBadRequest, "Bad Request", description}
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

// NotFound responds with error status code 404, Not Found.
func NotFound(c *gin.Context, description interface{}) {
	err := Error{http.StatusNotFound, "Not Found", description}
	c.JSON(http.StatusNotFound, gin.H{"error": err})
}

// InternalServerError responds with error status code 500, Internal Server Error.
func InternalServerError(c *gin.Context, description interface{}) {
	err := Error{http.StatusInternalServerError, "Internal Server Error", description}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err})
}
