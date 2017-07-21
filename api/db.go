package api

import (
	"database/sql"
	"os"

	raven "github.com/getsentry/raven-go"
	_ "github.com/go-sql-driver/mysql" // go mysql driver
)

// DBAuthInfo contains username, password and database name information.
// Used when making database connection.
type DBAuthInfo struct {
	User string // local mysql user
	Pass string // local mysql password
	Name string // database name
}

// DB is a global datastore for database
// connection credential information.
var DB DBAuthInfo

// SetDB updates DB with local connection information
// using enviroment variables.
func SetDB() {
	if DB.User = os.Getenv("DB_USER_V1"); DB.User == "" {
		panic("OS enviroment variable 'DB_USER_V1' is missing")
	}
	if DB.Pass = os.Getenv("DB_PASS_V1"); DB.Pass == "" {
		panic("OS enviroment variable 'DB_PASS_V1' is missing")
	}
	if DB.Name = os.Getenv("DB_NAME_V1"); DB.Name == "" {
		panic("OS enviroment variable 'DB_NAME_V1' is missing")
	}
}

// cred returns dataSourceName credentials string.
func cred() string {
	return DB.User + ":" + DB.Pass + "@tcp(127.0.0.1:3306)/" + DB.Name
}

// InitDB prepares database abstraction for later use.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", cred())
	if err != nil {
		raven.CaptureError(err, nil)
	}
	if err = db.Ping(); err != nil {
		raven.CaptureError(err, nil)
	}
	return db, err
}
