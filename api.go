package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gin-gonic/gin.v1"
)

type Cell struct {
	Id        int            `json:"id"`
	Accession sql.NullString `json:"accession"`
	Name      string         `json:"name"`
	Tissue    sql.NullString `json:"tissue"`
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

func getCell(c *gin.Context) {
	var (
		cell   Cell
		result gin.H
	)

	db := InitDb()
	defer db.Close()

	id := c.Param("id")
	row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_name = ?;", id)
	err := row.Scan(&cell.Id, &cell.Accession, &cell.Name, &cell.Tissue)
	if err != nil {
		row = db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
		err = row.Scan(&cell.Id, &cell.Accession, &cell.Name, &cell.Tissue)
	}
	if err != nil {
		row = db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.accession_id = ?;", id)
		err = row.Scan(&cell.Id, &cell.Accession, &cell.Name, &cell.Tissue)
	}

	if err != nil {
		result = gin.H{
			"category": "cell line",
			"count":    0,
			"data":     nil,
		}
		c.JSON(http.StatusNotFound, result)
	} else {
		result = gin.H{
			"category": "cell line",
			"count":    1,
			"data":     cell,
		}
		c.JSON(http.StatusOK, result)
	}
}

func main() {
	router := gin.Default()
	v1 := router.Group("v1")
	{
		v1.GET("/cell_lines/:id", getCell)
	}

	router.Run(":3000")
}
