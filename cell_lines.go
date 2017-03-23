package main

import (
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetCells handles GET requests for /cell_lines endpoint.
func GetCells(c *gin.Context) {
	getDataTypes(c, "List of all cell lines in pharmacodb", "select cell_id, cell_name from cells;")
}

// GetCellStats handles GET requests for /cell_lines/stats endpoint.
func GetCellStats(c *gin.Context) {
	getDataTypeStats(c, "Number of cell lines tested in each dataset",
		"select dataset_id, cell_lines from dataset_statistics;")
}

// GetCellIDs handles GET requests for /cell_lines/ids endpoint.
func GetCellIDs(c *gin.Context) {
	getDataTypeIDs(c, "List of all cell line IDs in pharmacodb", "select cell_id from cells;")
}

// GetCellNames handles GET requests for /cell_lines/names endpoint.
func GetCellNames(c *gin.Context) {
	getDataTypeNames(c, "List of all cell line names in pharmacodb", "select cell_name from cells;")
}

// GetCellByID handles GET requests for /cell_lines/ids/:id endpoint.
func GetCellByID(c *gin.Context) {
	var (
		cell      Cell
		syname    string
		synsource string
		syns      []Synonym
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	id := c.Param("id")
	queryStr := "select c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name, s.source_name, scn.cell_name from cells c inner join tissues t on t.tissue_id = c.tissue_id inner join source_cell_names scn on scn.cell_id = c.cell_id inner join sources s on s.source_id = scn.source_id where c.cell_id = ?"
	rows, err := db.Query(queryStr, id)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[string]bool)
	iter := 0
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Accession, &cell.Name, &cell.Tissue.ID, &cell.Tissue.Name, &synsource, &syname)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[syname] {
			for i, syn := range syns {
				if syn.Name == syname {
					syns[i].Datasets = append(syns[i].Datasets, synsource)
					break
				}
			}
		} else {
			var syn Synonym
			syn.Name = syname
			syn.Datasets = append(syn.Datasets, synsource)
			syns = append(syns, syn)
			exists[syname] = true
		}
		iter = 1
	}
	if iter == 0 {
		handleError(c, nil, http.StatusNotFound, fmt.Sprintf("Cell line with ID - %s - not found in pharmacodb", id))
		return
	}
	cell.Synonyms = syns

	c.IndentedJSON(http.StatusOK, gin.H{
		"type": "cell line",
		"data": cell,
	})
}
