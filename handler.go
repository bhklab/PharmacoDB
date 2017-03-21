package main

import (
	"database/sql"
	"os"

	"github.com/getsentry/raven-go"
	"gopkg.in/gin-gonic/gin.v1"
)

func init() {
	raven.SetDSN("https://24b828d4b8ea469da5b61941b6a3554a:d1de3fc962314598bcb3d04f010ce676@sentry.io/148972")
}

// ErrorHandler handles function error messages
func ErrorHandler(c *gin.Context, code int, message string) {
	c.IndentedJSON(code, gin.H{
		"error": gin.H{
			"status":  code,
			"message": message,
		},
	})
}

// InitDB prepares database abstraction for later use
func InitDB() *sql.DB {
	dbname := os.Getenv("pmdb_name")
	passwd := os.Getenv("mysql_passwd")
	cred := "root:" + passwd + "@tcp(127.0.0.1:3306)/" + dbname

	db, err := sql.Open("mysql", cred)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	return db
}
