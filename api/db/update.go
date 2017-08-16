package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // go mysql driver
)

// Group models the statistics data for each dataset.
type Group struct {
	ID         int
	Name       string
	Cell       int
	Tissue     int
	Drug       int
	Experiment int
}

// Groups is a collection of Group.
type Groups []Group

func main() {
	var groups Groups

	// Open db connection.
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleting old 'source_statistics' table ... ")

	// Delete statistics table.
	stmt, err := db.Prepare("DROP TABLE IF EXISTS source_statistics;")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("DONE\n")
	fmt.Printf("Creating new 'source_statistics' table ... ")

	// Create new statistics table.
	stmt, err = db.Prepare("CREATE TABLE source_statistics (dataset_id INT(11), dataset_name VARCHAR(255), cell_lines INT(11), tissues INT(11), drugs INT(11), experiments INT(11));")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("DONE\n")
	fmt.Printf("Getting current stats data for db update ... ")

	// Get stats current data.
	rows, err := db.Query("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.cell_id) FROM experiments e WHERE e.dataset_id = d.dataset_id) AS cell_lines, (SELECT COUNT(DISTINCT e.tissue_id) FROM experiments e WHERE e.dataset_id = d.dataset_id) AS tissues, (SELECT COUNT(DISTINCT e.drug_id) AS drugs FROM experiments e WHERE e.dataset_id = d.dataset_id) AS drugs, (SELECT COUNT(*) AS experiments FROM experiments e WHERE e.dataset_id = d.dataset_id) AS experiments FROM datasets d;")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var group Group
		err = rows.Scan(&group.ID, &group.Name, &group.Cell, &group.Tissue, &group.Drug, &group.Experiment)
		if err != nil {
			panic(err)
		}
		groups = append(groups, group)
	}

	fmt.Printf("DONE\n")
	fmt.Printf("Updating table with current data ... ")

	// Update statistics table with current data.
	for _, g := range groups {
		stmt, err := db.Prepare("INSERT INTO source_statistics (dataset_id, dataset_name, cell_lines, tissues, drugs, experiments) VALUES (?, ?, ?, ?, ?, ?);")
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec(g.ID, g.Name, g.Cell, g.Tissue, g.Drug, g.Experiment)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("DONE\n")
	fmt.Printf("Update finished successfully!\n")
}

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
// using environment variables.
func SetDB() {
	DB.User = os.Getenv("DB_USER_V1")
	DB.Pass = os.Getenv("DB_PASS_V1")
	DB.Name = os.Getenv("DB_NAME_V1")
}

// cred returns dataSourceName credentials string.
func cred() string {
	return DB.User + ":" + DB.Pass + "@tcp(127.0.0.1:3306)/" + DB.Name
}

// InitDB prepares database abstraction for later use.
func InitDB() (*sql.DB, error) {
	SetDB()
	db, err := sql.Open("mysql", cred())
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db, err
}
