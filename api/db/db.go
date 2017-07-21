package db

import (
	"database/sql"
	"os"

	"github.com/bhklab/pharmacodb/api"
	_ "github.com/go-sql-driver/mysql" // go mysql driver
)

// AuthInfo contains username, password and database name information.
// Used when making database connection.
type AuthInfo struct {
	User string // local mysql user
	Pass string // local mysql password
	Name string // database name
}

// DB is a global datastore for database
// connection credential information.
var DB AuthInfo

// Set updates DB with local connection information
// using enviroment variables.
func Set() {
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

// Cred returns dataSourceName string.
func Cred() string {
	return DB.User + ":" + DB.Pass + "@tcp(127.0.0.1:3306)/" + DB.Name
}

// Init prepares database abstraction for later use.
func Init() (*sql.DB, error) {
	db, err := sql.Open("mysql", Cred())
	if err != nil {
		api.LogPrivateError(err)
	}
	if err = db.Ping(); err != nil {
		api.LogPrivateError(err)
	}
	return db, err
}
