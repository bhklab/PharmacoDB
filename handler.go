package main

import (
	"database/sql"
	"log"
	"os"
)

// InitDb prepares database abstraction for later use
func InitDb() *sql.DB {
	dbname := os.Getenv("pharmacodb_api_dbname")
	passwd := os.Getenv("local_mysql_passwd")
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
