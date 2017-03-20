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

// GetCells handles GET requests for /cell_lines
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
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}

	rows, err := db.Query("select cell_id, cell_name from cells;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}
		cells = append(cells, cell)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"count":      len(cells),
		"cell_lines": cells,
	})
}

// GetCellStats handles GET requests for /cell_lines/stats
func GetCellStats(c *gin.Context) {
	var (
		stat  DatasetStat
		stats []DatasetStat
	)

	db := InitDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}

	rows, err := db.Query("select dataset_id, cell_lines from dataset_statistics;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&stat.Dataset, &stat.Count)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}
		stats = append(stats, stat)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"statistics": stats,
	})
}
