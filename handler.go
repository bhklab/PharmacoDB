package main

import (
	"database/sql"
	"log"
	"os"

	"gopkg.in/gin-gonic/gin.v1"
)

// Synonym is a collection of synonyms for a resource
// Resource is one of: cell_line, tissue, drug
type Synonym struct {
	Name    string   `json:"name"`
	Sources []string `json:"sources"`
}

// InitDb prepares database abstraction for later use
func InitDb() *sql.DB {
	dbname := os.Getenv("pmdb_name")
	passwd := os.Getenv("mysql_passwd")
	cred := "root:" + passwd + "@tcp(127.0.0.1:3306)/" + dbname

	db, err := sql.Open("mysql", cred)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// ErrHandler handles all error messages
// Does not handle middleware errors (Eg. bad routes/endpoints)
func ErrHandler(c *gin.Context, status int, message string) {
	c.IndentedJSON(status, gin.H{
		"error": gin.H{
			"status":  status,
			"message": message,
		},
	})
}
