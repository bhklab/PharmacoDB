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
func IndexCell(c *gin.Context) {
	var (
		cell  Cell
		cells []Cell
		total int
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	all := c.DefaultQuery("all", "false")
	if all == "true" {
		rows, er := db.Query("SELECT cell_id, accession_id, cell_name FROM cells;")
		defer rows.Close()
		if er != nil {
			handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
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
			"total":       len(cells),
			"data":        cells,
			"description": "List of all cell lines in PharmacoDB",
		})
		return
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "30") // Need to fix the off-by-1 results

	// Add handler for off-limit queries here
	// Need testing for pages outside data limit

	query := "SELECT SQL_CALC_FOUND_ROWS cell_id, accession_id, cell_name FROM cells limit " + page + "," + limit + ";"
	rows, err := db.Query(query)
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
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":        cells,
		"description": "List of cell lines in PharmacoDB",
		"page":        page,
		"limit":       limit,
		"total":       total,
	})
}
