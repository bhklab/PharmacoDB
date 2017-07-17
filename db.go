package main

import (
	"database/sql"
	"os"
)

// DBAuthInfo contains username, password and database name information.
// Used when making database connection.
type DBAuthInfo struct {
	User string // local mysql user
	Pass string // password
	Name string // local mysql database name
}

// DB is a global datastore for database
// connection credential information.
var DB DBAuthInfo

// SetDB sets DB with local connection data
// using environment variables.
func SetDB() {
	if DB.User = os.Getenv("MYSQL_USER_V1"); DB.User == "" {
		panic("OS enviroment variable 'MYSQL_USER_V1' is missing")
	}
	if DB.Pass = os.Getenv("MYSQL_PASS_V1"); DB.Pass == "" {
		panic("OS enviroment variable 'MYSQL_PASS_V1' is missing")
	}
	if DB.Name = os.Getenv("MYSQL_NAME_V1"); DB.Name == "" {
		panic("OS enviroment variable 'MYSQL_NAME_V1' is missing")
	}
}

// DBCred returns a credentials string for
// making a database connection.
func DBCred() string {
	return DB.User + ":" + DB.Pass + "@tcp(127.0.0.1:3306)/" + DB.Name
}

// InitDB prepares database abstraction for later use.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", DBCred())
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
	}
	err = db.Ping()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
	}
	return db, err
}
