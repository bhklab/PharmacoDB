package main

import (
	"database/sql"
	"os"

	raven "github.com/getsentry/raven-go"
	"gopkg.in/gin-gonic/gin.v1"
)

// Set Sentry DSN for internal error logging.
func init() {
	raven.SetDSN("https://24b828d4b8ea469da5b61941b6a3554a:d1de3fc962314598bcb3d04f010ce676@sentry.io/148972")
}

// Prepare database abstraction for later use.
func initDB() (*sql.DB, error) {
	cred := os.Getenv("mysql_user") + ":" + os.Getenv("mysql_passwd") + "@tcp(127.0.0.1:3306)/" + os.Getenv("pmdb_name")
	db, err := sql.Open("mysql", cred)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	err = db.Ping()
	if err != nil {
		raven.CaptureError(err, nil)
	}
	return db, err
}

// Handle request error messages (all except no route match errors).
func handleError(c *gin.Context, err error, code int, message string) {
	if err != nil {
		raven.CaptureError(err, nil)
	}
	c.IndentedJSON(code, gin.H{
		"error": gin.H{
			"status":  code,
			"message": message,
		},
	})
}
