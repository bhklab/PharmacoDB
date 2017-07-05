package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Experiment is an experiment datatype.
type Experiment struct {
	ID      int     `json:"experiment_id"`
	Cell    Cell    `json:"cell_line"`
	Tissue  Tissue  `json:"tissue"`
	Drug    Drug    `json:"drug"`
	Dataset Dataset `json:"dataset"`
}

// IndexExperiment ...
func IndexExperiment(c *gin.Context) {
	var (
		experiment  Experiment
		experiments []Experiment
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	curPage := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("per_page", "10")

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
	SQL1 := "SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, "
	SQL2 := "d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e "
	SQL3 := "JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id "
	SQL4 := "JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id"
	selectSQL := SQL1 + SQL2 + SQL3 + SQL4
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		experiments = append(experiments, experiment)
	}
	row := db.QueryRow("SELECT COUNT(*) FROM experiments;")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/experiments", page, total, limit)

	c.IndentedJSON(http.StatusOK, gin.H{
		"data":        experiments,
		"total":       total,
		"description": "List of all experiments in PharmacoDB",
	})
}
