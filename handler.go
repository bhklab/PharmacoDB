package main

import (
	"database/sql"
	"os"

	"github.com/getsentry/raven-go"
	"gopkg.in/gin-gonic/gin.v1"
)

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
