package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IndexExperiment returns a list of all experiments currently in database.
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

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Set max limit per_page to 100
	if limit > 100 {
		limit = 100
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

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        experiments,
			"total":       total,
			"description": "List of all experiments in PharmacoDB",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        experiments,
			"total":       total,
			"description": "List of all experiments in PharmacoDB",
		})
	}
}

// ShowExperiment returns dose response data for a specific experiment id.
func ShowExperiment(c *gin.Context) {
	var (
		experiment   Experiment
		doseResponse DoseResponse
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")
	SQL1 := "SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, "
	SQL2 := "d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e "
	SQL3 := "JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id "
	SQL4 := "JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.experiment_id = ?;"
	selectSQL := SQL1 + SQL2 + SQL3 + SQL4
	row := db.QueryRow(selectSQL, id)
	err = row.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Sprintf("Experiment with ID:%s not found in database", id)
			handleError(c, nil, http.StatusNotFound, errMessage)
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	rows, err := db.Query("SELECT dose, response FROM dose_responses WHERE experiment_id = ?;", id)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&doseResponse.Dose, &doseResponse.Response)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		experiment.DR = append(experiment.DR, doseResponse)
	}

	desc := fmt.Sprintf("Dose/Response data for experiment with ID:%s", id)
	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        experiment,
			"description": desc,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        experiment,
			"description": desc,
		})
	}
}
