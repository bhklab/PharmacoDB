package main

import (
	// "fmt"
	"os"
	"log"
	// "net/http"
	"database/sql"

	// "gopkg.in/gin-gonic/gin.v1"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbname := os.Getenv("pharmacodb_api_dbname")
	passwd := os.Getenv("local_mysql_passwd")

	cred := "root:" + passwd + "@tcp(127.0.0.1:3306)/" + dbname
	
	// prepare database abstraction for later use
	db, err := sql.Open("mysql", cred)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// check that a network connection can be established and login
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

}
