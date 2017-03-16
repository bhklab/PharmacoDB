package main

import (
	"log"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/guregu/null.v3"
)

// Cell is a cell line
type Cell struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Accession null.String `json:"accession"`
	Tissue    string      `json:"tissue,omitempty"`
}

// GetCLines handles GET requests for cell lines
func GetCLines(c *gin.Context) {
	var (
		cell  Cell
		cells []Cell
	)

	db := InitDb()
	defer db.Close()

	rows, err := db.Query("select cell_id, accession_id, cell_name from cells;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Accession, &cell.Name)
		if err != nil {
			log.Fatal(err)
		}
		cells = append(cells, cell)
	}
	defer rows.Close()

	result := gin.H{
		"category":    "cell_lines",
		"description": "list of all cell lines in pharmacodb",
		"count":       len(cells),
		"data":        cells,
	}
	c.IndentedJSON(http.StatusOK, result)
}

// GetCLineByID handles GET request for a cell line using ID
func GetCLineByID(c *gin.Context) {
	var cell Cell

	db := InitDb()
	defer db.Close()

	id := c.Param("id")
	row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
	err := row.Scan(&cell.ID, &cell.Accession, &cell.Name, &cell.Tissue)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"status":  http.StatusNotFound,
				"message": "cell line with id - " + id + " - not found in database",
			},
		})
	} else {
		c.IndentedJSON(http.StatusOK, cell)
	}
}
