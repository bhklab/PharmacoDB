package main

import (
	"log"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetCLines : GET request handler for cell lines
func GetCLines(c *gin.Context) {
	var (
		cell  Cells
		cells []Cells
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
		"category": "cell line",
		"count":    len(cells),
		"data":     cells,
	}
	c.JSON(http.StatusOK, result)
}
