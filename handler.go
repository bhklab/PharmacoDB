package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/getsentry/raven-go"
	"gopkg.in/gin-gonic/gin.v1"
)

// DatasetStat contains the number of a resource tested in a dataset
type DatasetStat struct {
	Dataset int `json:"dataset"`
	Count   int `json:"count"`
}

func init() {
	raven.SetDSN("https://24b828d4b8ea469da5b61941b6a3554a:d1de3fc962314598bcb3d04f010ce676@sentry.io/148972")
}

// ErrorHandler handles func error messages
func ErrorHandler(c *gin.Context, code int, message string) {
	c.IndentedJSON(code, gin.H{
		"error": gin.H{
			"status":  code,
			"message": message,
		},
	})
}

// InitDb prepares database abstraction for later use
func InitDb() *sql.DB {
	dbname := os.Getenv("pmdb_name")
	passwd := os.Getenv("mysql_passwd")
	cred := "root:" + passwd + "@tcp(127.0.0.1:3306)/" + dbname

	db, err := sql.Open("mysql", cred)
	if err != nil {
		raven.CaptureError(err, nil)
	}

	return db
}

// GetDataTypes handles request for /datatypes
func GetDataTypes(c *gin.Context) {
	data := [5]string{"cell_lines", "tissues", "drugs", "datasets", "experiments"}
	c.IndentedJSON(http.StatusOK, gin.H{
		"data": data,
	})
}
