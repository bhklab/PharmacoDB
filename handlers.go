package main

import (
	"database/sql"
	"net/http"
	"os"

	raven "github.com/getsentry/raven-go"
	"gopkg.in/gin-gonic/gin.v1"
)

// Set Sentry DSN for internal error logging.
func init() {
	raven.SetDSN("https://24b828d4b8ea469da5b61941b6a3554a:d1de3fc962314598bcb3d04f010ce676@sentry.io/148972")
}

// Prepare database abstraction for later use.
func initDB() (*sql.DB, error) {
	cred := os.Getenv("mysql_user") + ":" + os.Getenv("mysql_passwd") + "@tcp(127.0.0.1:3306)/" + os.Getenv("pmdb_name")
	db, err := sql.Open("mysql", cred)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	err = db.Ping()
	if err != nil {
		raven.CaptureError(err, nil)
	}
	return db, err
}

// Handle request error messages (all except no route match errors).
func handleError(c *gin.Context, err error, code int, message string) {
	if err != nil {
		raven.CaptureError(err, nil)
	}
	c.IndentedJSON(code, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

// getDataTypes is an abstract GET request handler for /{datatype} endpoints.
// Endpoints: /cell_lines, /tissues, /drugs, /datasets
func getDataTypes(c *gin.Context, desc string, queryStr string) {
	var (
		item  DataTypeReduced
		items []DataTypeReduced
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	rows, err := db.Query(queryStr)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&item.ID, &item.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		items = append(items, item)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": desc,
		"count":       len(items),
		"data":        items,
	})
}

// getDataTypeStats is an abstract GET request handler for /{datatype}/stats endpoints.
// Endpoints: /cell_lines/stats, /tissues/stats, /drugs/stats
func getDataTypeStats(c *gin.Context, desc string, queryStr string) {
	var (
		stat  DatasetStat
		stats []DatasetStat
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	rows, err := db.Query(queryStr)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&stat.Dataset, &stat.Count)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		stats = append(stats, stat)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": desc,
		"data":        stats,
	})
}

// getDataTypeIDs is an abstract GET request handler for /{datatypes}/ids endpoints.
// Endpoints: /cell_lines/ids, /tissues/ids, /drugs/ids, /datasets/ids
func getDataTypeIDs(c *gin.Context, desc string, queryStr string) {
	var (
		id  int
		ids []int
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	rows, err := db.Query(queryStr)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		ids = append(ids, id)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": desc,
		"data":        ids,
	})
}

// getDataTypeNames is an abstract GET request handler for /{datatypes}/names endpoints.
// Endpoints: /cell_lines/names, /tissues/names, /drugs/names, /datasets/names
func getDataTypeNames(c *gin.Context, desc string, queryStr string) {
	var (
		name  string
		names []string
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	rows, err := db.Query(queryStr)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		names = append(names, name)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": desc,
		"data":        names,
	})
}
