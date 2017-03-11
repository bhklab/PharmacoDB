package main

import (
	"database/sql"
	// "encoding/json"
	// "fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gin-gonic/gin.v1"
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

	type Cell struct {
		Cell_Id      int            `json:"id"`
		Accession_Id sql.NullString `json:"accession"`
		Cell_Name    string         `json:"name"`
	}

	router := gin.Default()

	v1 := router.Group("v1")
	{
		// handle GET request
		// Endpoints: /cell_lines
		v1.GET("/cell_lines", func(c *gin.Context) {
			var (
				cell  Cell
				cells []Cell
			)
			rows, err := db.Query("select cell_id, accession_id, cell_name from cells;")
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				err = rows.Scan(&cell.Cell_Id, &cell.Accession_Id, &cell.Cell_Name)
				cells = append(cells, cell)
				if err != nil {
					log.Fatal(err)
				}
			}
			defer rows.Close()

			result := gin.H{
				"category": "cell line",
				"count":    len(cells),
				"data":     cells,
			}
			c.JSON(http.StatusOK, result)
		})

		// handle GET request (for all three endpoints below)
		// Endpoints: /cell_lines/:id, /cell_lines/:name, /cell_lines/:accession
		v1.GET("/cell_lines/:id", func(c *gin.Context) {
			var (
				cell   Cell
				result gin.H
			)
			id := c.Param("id")
			row := db.QueryRow("select cell_id, accession_id, cell_name from cells where cell_name = ?;", id)
			err = row.Scan(&cell.Cell_Id, &cell.Accession_Id, &cell.Cell_Name)
			if err != nil {
				row := db.QueryRow("select cell_id, accession_id, cell_name from cells where cell_id = ?;", id)
				err = row.Scan(&cell.Cell_Id, &cell.Accession_Id, &cell.Cell_Name)
				if err != nil {
					row := db.QueryRow("select cell_id, accession_id, cell_name from cells where accession_id = ?;", id)
					err = row.Scan(&cell.Cell_Id, &cell.Accession_Id, &cell.Cell_Name)
					if err != nil {
						result = gin.H{
							"category": "cell line",
							"count":    0,
							"data":     nil,
						}
					} else {
						result = gin.H{
							"category": "cell line",
							"count":    1,
							"data":     cell,
						}
					}
				} else {
					result = gin.H{
						"category": "cell line",
						"count":    1,
						"data":     cell,
					}
				}
			} else {
				result = gin.H{
					"category": "cell line",
					"count":    1,
					"data":     cell,
				}
			}
			c.JSON(http.StatusOK, result)
		})
	}
	gin.SetMode(gin.ReleaseMode)
	router.Run(":3000")
}
