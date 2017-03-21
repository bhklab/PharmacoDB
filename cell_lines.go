package main

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetCells handles GET requests for /cell_lines endpoint.
func GetCells(c *gin.Context) {
	var (
		cell  CellReduced
		cells []CellReduced
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	rows, err := db.Query("select cell_id, cell_name from cells;")
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		cells = append(cells, cell)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"count":      len(cells),
		"cell_lines": cells,
	})
}
