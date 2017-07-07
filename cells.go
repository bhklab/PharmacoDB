package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Cell is a cell_line datatype.
type Cell struct {
	ID     int     `json:"id"`
	ACC    *string `json:"accession_id,omitempty"`
	Name   string  `json:"name"`
	Tissue *Tissue `json:"tissue,omitempty"`
}

// IndexCell returns a list of all cell lines currently in database.
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

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if isTrue, _ := strconv.ParseBool(c.DefaultQuery("all", "true")); isTrue {
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
		if shouldIndent {
			c.IndentedJSON(http.StatusOK, gin.H{
				"data":        cells,
				"total":       len(cells),
				"description": "List of all cell lines in PharmacoDB",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":        cells,
				"total":       len(cells),
				"description": "List of all cell lines in PharmacoDB",
			})
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	s := (page - 1) * limit
	selectSQL := "SELECT cell_id, accession_id, cell_name FROM cells"
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
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
	row := db.QueryRow("SELECT COUNT(*) FROM cells;")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/cell_lines", page, total, limit)

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        cells,
			"total":       total,
			"description": "List of all cell lines in PharmacoDB",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        cells,
			"total":       total,
			"description": "List of all cell lines in PharmacoDB",
		})
	}
}

// ShowCell ...
func ShowCell(c *gin.Context) {

}
