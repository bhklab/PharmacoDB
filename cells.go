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
	ACC    *string `json:"accession_id"`
	Name   string  `json:"name"`
	Tissue *Tissue `json:"tissue,omitempty"`
}

// IndexCell returns a list of all cell lines currently in database.
// To get all data in one request: /cell_lines?all=true
// To get paginated data: /cell_lines?page=1&per_page=30 (links in header)
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

	all := c.DefaultQuery("all", "false")
	if isTrue, _ := strconv.ParseBool(all); isTrue {
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
			"data":        cells,
			"total":       len(cells),
			"description": "List of all cell lines in PharmacoDB",
		})
		return
	}

	curPage := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("per_page", "30")

	page, err := strconv.Atoi(curPage)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	limit, err := strconv.Atoi(perPage)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS cell_id, accession_id, cell_name FROM cells limit %d,%d;", s, limit)
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
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/cell_lines", page, total, limit)

	c.IndentedJSON(http.StatusOK, gin.H{
		"data":        cells,
		"total":       total,
		"description": "List of all cell lines in PharmacoDB",
	})
}
