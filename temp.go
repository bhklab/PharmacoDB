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

var (
	cell   Cell
	result gin.H
)

func getCell(c *gin.Context) {
	id := c.Param("id")
	row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
	err = row.Scan(&cell.Id, &cell.Accession, &cell.Name, &cell.Tissue)
	if err != nil {
		row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
		err = row.Scan(&cell.Id, &cell.Accession, &cell.Name, &cell.Tissue)
		if err != nil {
			row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
			err = row.Scan(&cell.Id, &cell.Accession, &cell.Name, &cell.Tissue)
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
}

func getCells(c *gin.Context) {
	var (
		cell   Cell
		result gin.H
	)
	id := c.Param("id")
	row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
	err = row.Scan(&cell.Cell_Id, &cell.Accession_Id, &cell.Cell_Name)
	if err != nil {
		row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
		err = row.Scan(&cell.Cell_Id, &cell.Accession_Id, &cell.Cell_Name)
		if err != nil {
			row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
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
}

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

// GET cell line request handler (param: id)
func CellByID(c *gin.Context) {
	var cell Cell

	db := InitDb()
	defer db.Close()

	id := c.Param("id")
	row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
	err := row.Scan(&cell.ID, &cell.Accession, &cell.Name, &cell.Tissue)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"status":  http.StatusNotFound,
				"message": fmt.Sprintf("cell line with ID - %s - not found in database", id),
			},
		})
	} else {
		c.JSON(http.StatusOK, cell)
	}
}

// GET cell line request handler (param: accession)
func CellByAcc(c *gin.Context) {
	var cell Cell

	db := InitDb()
	defer db.Close()

	acc := c.Param("acc")
	row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.accession_id = ?;", acc)
	err := row.Scan(&cell.ID, &cell.Accession, &cell.Name, &cell.Tissue)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"status":  http.StatusNotFound,
				"message": fmt.Sprintf("cell line with accession - %s - not found in database", acc),
			},
		})
	} else {
		c.JSON(http.StatusOK, cell)
	}
}
