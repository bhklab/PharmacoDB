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
		v1.GET("/cell_lines", getCells)
		v1.GET("/cell_lines/names", getCellsNames)
		v1.GET("/cell_lines/names/:name", getCellByName)
		v1.GET("/cell_lines/names/:name/synonyms", getCellSynByName)
		v1.GET("/cell_lines/names/:name/drugs", getCellDrugsByName)
		v1.GET("/cell_lines/names/:name/drugs_stat", getCellDrugStatByName)
		v1.GET("/cell_lines/ids", getCellsIDs)
		v1.GET("/cell_lines/ids/:id", getCellByID)
		v1.GET("/cell_lines/ids/:id/synonyms", getCellSynByID)
		v1.GET("/cell_lines/ids/:id/drugs", getCellDrugsByID)
		v1.GET("/cell_lines/ids/:id/drugs_stat", getCellDrugStatByID)
	}

	router.Run(":3000")
}

// prepare database abstraction for later use
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
