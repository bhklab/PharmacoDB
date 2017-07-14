package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CellCountPerTissue returns the cell line count per tissue for all tissues.
func CellCountPerTissue(c *gin.Context) {
	type TissueCellCount struct {
		Tissue Tissue `json:"tissue"`
		Count  int    `json:"cell-line-count"`
	}
	var (
		tcc  TissueCellCount
		tccs []TissueCellCount
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	q1 := "SELECT t.tissue_id, t.tissue_name, (SELECT COUNT(DISTINCT c.cell_id) "
	q2 := "FROM cells c WHERE c.tissue_id = t.tissue_id) as cell_count FROM tissues t;"
	query := q1 + q2
	rows, er := db.Query(query)
	defer rows.Close()
	if er != nil {
		handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&tcc.Tissue.ID, &tcc.Tissue.Name, &tcc.Count)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		tccs = append(tccs, tcc)
	}

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":          tccs,
			"total-tissues": len(tccs),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":          tccs,
			"total-tissues": len(tccs),
		})
	}
}

// DrugCountPerDataset returns the number of drugs tested in each dataset.
func DrugCountPerDataset(c *gin.Context) {
	type DatasetDrugCount struct {
		Dataset Dataset `json:"dataset"`
		Count   int     `json:"drug-count"`
	}
	var (
		ddc  DatasetDrugCount
		ddcs []DatasetDrugCount
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	q1 := "SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.drug_id) "
	q2 := "FROM experiments e WHERE e.dataset_id = d.dataset_id) as drug_count FROM datasets d;"
	query := q1 + q2
	rows, er := db.Query(query)
	defer rows.Close()
	if er != nil {
		handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&ddc.Dataset.ID, &ddc.Dataset.Name, &ddc.Count)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		ddcs = append(ddcs, ddc)
	}

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":           ddcs,
			"total-datasets": len(ddcs),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":           ddcs,
			"total-datasets": len(ddcs),
		})
	}
}

// CellDrugsPerDataset returns the number of drugs tested against a particular cell line per dataset.
func CellDrugsPerDataset(c *gin.Context) {
	type DatasetDrugCount struct {
		Dataset Dataset `json:"dataset"`
		Count   int     `json:"drug-count"`
	}
	var (
		ddc  DatasetDrugCount
		ddcs []DatasetDrugCount
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")

	q1 := "SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.drug_id) "
	q2 := "FROM experiments e WHERE e.dataset_id = d.dataset_id AND e.cell_id = ?) as drug_count FROM datasets d;"
	query := q1 + q2
	rows, er := db.Query(query, id)
	defer rows.Close()
	if er != nil {
		handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&ddc.Dataset.ID, &ddc.Dataset.Name, &ddc.Count)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		ddcs = append(ddcs, ddc)
	}

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":           ddcs,
			"total-datasets": len(ddcs),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":           ddcs,
			"total-datasets": len(ddcs),
		})
	}
}

// CellDrugExperiments returns all dose/response data for a cell line and drug combination.
func CellDrugExperiments(c *gin.Context) {
	var (
		experiment   Experiment
		doseResponse DoseResponse
		experiments  []Experiment
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	cellID := c.Param("cell_id")
	drugID := c.Param("drug_id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, "
	SQL2 := "d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e "
	SQL3 := "JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id "
	SQL4 := "JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE "
	var SQL5 string
	if searchByName(searchType) {
		SQL5 = "c.cell_name LIKE ? AND d.drug_name LIKE ?;"
	} else {
		SQL5 = "e.cell_id = ? AND e.drug_id = ?;"
	}
	selectSQL := SQL1 + SQL2 + SQL3 + SQL4 + SQL5
	rrows, _ := db.Query(selectSQL, cellID, drugID)
	for rrows.Next() {
		err = rrows.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		query := "SELECT dose, response FROM dose_responses WHERE experiment_id = ?;"
		rows, err := db.Query(query, experiment.ID)
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
		experiments = append(experiments, experiment)
	}

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data": experiments,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": experiments,
		})
	}
}

// CellDatasetExperiments returns all experiments where a cell line has been tested against a drug,
// in a specified dataset.
func CellDatasetExperiments(c *gin.Context) {
	var (
		experiment   Experiment
		doseResponse DoseResponse
		experiments  []Experiment
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	cellID := c.Param("cell_id")
	datasetID := c.Param("dataset_id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, "
	SQL2 := "d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e "
	SQL3 := "JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id "
	SQL4 := "JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE "
	var SQL5 string
	if searchByName(searchType) {
		SQL5 = "c.cell_name LIKE ? AND da.dataset_name LIKE ?;"
	} else {
		SQL5 = "e.cell_id = ? AND e.dataset_id = ?;"
	}
	selectSQL := SQL1 + SQL2 + SQL3 + SQL4 + SQL5
	rrows, _ := db.Query(selectSQL, cellID, datasetID)
	for rrows.Next() {
		err = rrows.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		query := "SELECT dose, response FROM dose_responses WHERE experiment_id = ?;"
		rows, err := db.Query(query, experiment.ID)
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
		experiments = append(experiments, experiment)
	}

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data": experiments,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": experiments,
		})
	}
}
