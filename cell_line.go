package main

import (
	"net/http"

	"github.com/getsentry/raven-go"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/guregu/null.v3"
)

// CellReduced is a cell line with only two attributes
type CellReduced struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Cell is a cell line datatype
type Cell struct {
	ID        int         `json:"id"`
	Accession null.String `json:"accession"`
	Name      string      `json:"name"`
	Tissue    Tissue      `json:"tissue"`
}

// GetCells handles GET requests for /cell_lines endpoint
func GetCells(c *gin.Context) {
	var (
		cell  CellReduced
		cells []CellReduced
	)

	db := InitDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
		c.Abort()
		return
	}

	rows, err := db.Query("select cell_id, cell_name from cells;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
			c.Abort()
			return
		}
		cells = append(cells, cell)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"count": len(cells),
		"data":  cells,
	})
}
