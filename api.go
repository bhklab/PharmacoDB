package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gin-gonic/gin.v1"
)

type Cell struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Accession string `json:"accession"`
	Tissue    string `json:"tissue"`
}

func main() {
	router := gin.Default()
	v1 := router.Group("v1")
	{
		v1.GET("/cell_lines/name/:name", CellByName)
	}

	router.Run(":3000")
}

// create a database connection
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

// GET cell line request handler (param: name)
func CellByName(c *gin.Context) {
	var cell Cell

	db := InitDb()
	defer db.Close()

	name := c.Param("name")
	row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_name = ?;", name)
	err := row.Scan(&cell.ID, &cell.Accession, &cell.Name, &cell.Tissue)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"status":  http.StatusNotFound,
				"message": fmt.Sprintf("cell line - %s - not found in database", name),
			},
		})
	} else {
		c.JSON(http.StatusOK, cell)
	}
}
