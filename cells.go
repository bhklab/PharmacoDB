package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cell is a cell_line datatype.
type Cell struct {
	ID   int     `json:"id"`
	ACC  *string `json:"accession_id"`
	Name string  `json:"name"`
}

// IndexCell ...
// Add ?page=123 && ?limit=123 fields handling feature to method
// For example: http://localhost:4200/v1/cell_lines?page=1&limit=30
func IndexCell(c *gin.Context) {
	var (
		cell  Cell
		cells []Cell
	)
	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	rows, err := db.Query("SELECT SQL_CALC_FOUND_ROWS cell_id, accession_id, cell_name FROM cells limit 1,30;")
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	total, err := db.Query(" SELECT FOUND_ROWS();")
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		cells = append(cells, cell)
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"count":       len(cells),
		"data":        cells,
		"description": "List of all cell lines in PharmacoDB",
	})
}
