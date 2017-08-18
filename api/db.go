package api

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql" // go mysql driver
)

// DBAuthInfo contains local mysql username, password, and database name.
type DBAuthInfo struct {
	User string // local mysql user
	Pass string // local mysql password
	Name string // database name
	Host string // mysql host
}

// DB is a global datastore for database connection information.
var DB DBAuthInfo

// SetDB updates DB using environment settings.
func SetDB(version string) {
	if DB.User = os.Getenv("DB_USER_V" + version); DB.User == "" {
		panic("Missing environment variable: DB_USER_V" + version)
	}
	if DB.Pass = os.Getenv("DB_PASS_V" + version); DB.Pass == "" {
		panic("Missing environment variable: DB_PASS_V" + version)
	}
	if DB.Name = os.Getenv("DB_NAME_V" + version); DB.Name == "" {
		panic("Missing environment variable: DB_NAME_V" + version)
	}
	if DB.Host = os.Getenv("DB_HOST_V" + version); DB.Host == "" {
		panic("Missing environment variable: DB_HOST_V" + version)
	}
}

// Database creates a new database connection.
func Database() (*sql.DB, error) {
	cred := DB.User + ":" + DB.Pass + "@tcp(" + DB.Host + ":3306)/" + DB.Name
	db, _ := sql.Open("mysql", cred)
	err := db.Ping()
	if err != nil {
		LogSentry(err)
	}
	return db, err
}
