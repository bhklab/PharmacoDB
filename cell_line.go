package main

import (
	"fmt"
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
	ID        int           `json:"id"`
	Accession null.String   `json:"accession"`
	Name      string        `json:"name"`
	Tissue    TissueReduced `json:"tissue"`
	Syn       []Synonym     `json:"synonyms"`
}

// GetCells handles GET requests for /cell_lines
func GetCells(c *gin.Context) {
	var (
		cell  CellReduced
		cells []CellReduced
	)

	db := InitDB()
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

	db := InitDB()
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

// GetCellIDs handles GET requests for /cell_lines/ids
func GetCellIDs(c *gin.Context) {
	db := InitDB()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}

	var (
		ID  string
		IDs []string
	)
	rows, dberr := db.Query("select cell_id from cells;")
	if dberr != nil {
		raven.CaptureError(dberr, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&ID)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}
		IDs = append(IDs, ID)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, IDs)
}

// GetCellByID handles GET requests for /cell_lines/ids/:id
func GetCellByID(c *gin.Context) {
	var (
		cell    Cell
		scname  ReSynonym
		scnames []ReSynonym
		scn     Synonym
		scns    []Synonym
	)

	db := InitDB()
	defer db.Close()

	dberr := db.Ping()
	if dberr != nil {
		raven.CaptureError(dberr, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}

	id := c.Param("id")

	rows, err := db.Query("select c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name, s.source_name, scn.cell_name from cells c inner join tissues t on t.tissue_id = c.tissue_id inner join source_cell_names scn on scn.cell_id = c.cell_id inner join sources s on s.source_id = scn.source_id where c.cell_id = ?", id)
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}

	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Accession, &cell.Name, &cell.Tissue.ID, &cell.Tissue.Name, &scname.Dataset, &scname.Name)
		if err != nil {
			ErrorHandler(c, http.StatusNotFound, fmt.Sprintf("Cell line with ID - %s - not found in database.", id))
			c.Abort()
			return
		}
		scnames = append(scnames, scname)
	}
	defer rows.Close()

	scnameHash := make(map[string]bool)
	for _, syn := range scnames {
		if scnameHash[syn.Name] {
			for _, b := range scns {
				if b.Name == syn.Name {
					b.Datasets = append(b.Datasets, syn.Dataset)
					fmt.Println(syn.Dataset)
					// something wrong is happening here (to be continued ...)
				}
			}
		} else {
			scnameHash[syn.Name] = true
			scn.Name = syn.Name
			var emptystr []string
			emptystr = append(emptystr, syn.Dataset)
			scn.Datasets = emptystr
			scns = append(scns, scn)
		}
	}

	cell.Syn = scns

	c.IndentedJSON(http.StatusOK, cell)
}
