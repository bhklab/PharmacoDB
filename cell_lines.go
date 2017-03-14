package main

import (
	"log"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/guregu/null.v3"
)

// Stat is a dataset-to-count statistics
// (used for plotting in webapp)
type Stat struct {
	Dataset string `json:"dataset"`
	Count   int    `json:"count"`
}

// Cell is a cell line (attr: 4)
type Cell struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Accession null.String `json:"accession"`
	Tissue    string      `json:"tissue"`
}

// Cells is a cell line (attr: 3)
type Cells struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Accession null.String `json:"accession"`
}

// GetCLines handles GET requests for cell lines
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
		"category": "cell_lines",
		"count":    len(cells),
		"data":     cells,
	}
	c.JSON(http.StatusOK, result)
}

// GetCLinesStats handles GET requests for cell lines statistics
func GetCLinesStats(c *gin.Context) {
	var (
		stat  Stat
		stats []Stat
	)

	db := InitDb()
	defer db.Close()

	rows, err := db.Query("select dataset_name, cell_lines from dataset_statistics;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&stat.Dataset, &stat.Count)
		if err != nil {
			log.Fatal(err)
		}
		stats = append(stats, stat)
	}
	defer rows.Close()

}
