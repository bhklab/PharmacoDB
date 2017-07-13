package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Drug is a drug datatype.
type Drug struct {
	ID   int       `json:"id"`
	Name string    `json:"name"`
	SYNS []Synonym `json:"synonyms,omitempty"`
}

// IndexDrug returns a list of all drugs currently in database.
func IndexDrug(c *gin.Context) {
	var (
		drug  Drug
		drugs []Drug
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if isTrue, _ := strconv.ParseBool(c.DefaultQuery("all", "false")); isTrue {
		rows, er := db.Query("SELECT drug_id, drug_name FROM drugs;")
		defer rows.Close()
		if er != nil {
			handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		for rows.Next() {
			err = rows.Scan(&drug.ID, &drug.Name)
			if err != nil {
				handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			drugs = append(drugs, drug)
		}
		if shouldIndent {
			c.IndentedJSON(http.StatusOK, gin.H{
				"data":        drugs,
				"total":       len(drugs),
				"description": "List of all drugs in PharmacoDB",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":        drugs,
				"total":       len(drugs),
				"description": "List of all drugs in PharmacoDB",
			})
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	s := (page - 1) * limit
	selectSQL := "SELECT drug_id, drug_name FROM drugs"
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		drugs = append(drugs, drug)
	}
	row := db.QueryRow("SELECT COUNT(*) FROM drugs;")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/drugs", page, total, limit)

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        drugs,
			"total":       total,
			"description": "List of all drugs in PharmacoDB",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        drugs,
			"total":       total,
			"description": "List of all drugs in PharmacoDB",
		})
	}
}

// ShowDrug returns a drug using ID or Name.
func ShowDrug(c *gin.Context) {
	var (
		drug     Drug
		synonym  Synonym
		synonyms []Synonym
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT drug_id, drug_name FROM drugs WHERE "
	var SQL2 string
	if searchByName(searchType) {
		SQL2 = "drug_name LIKE ?;"
	} else {
		SQL2 = "drug_id LIKE ?;"
	}
	SQL := SQL1 + SQL2
	row := db.QueryRow(SQL, id)
	err = row.Scan(&drug.ID, &drug.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			message := fmt.Sprintf("Drug with ID:%s not found in database", id)
			handleError(c, nil, http.StatusNotFound, message)
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	q1 := "SELECT s.drug_name, d.dataset_name FROM source_drug_names s "
	q2 := "JOIN datasets d ON d.dataset_id = s.source_id WHERE s.drug_id = ?;"
	query := q1 + q2
	rows, err := db.Query(query, drug.ID)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	var (
		synonymName string
		datasetName string
	)
	exists := make(map[string]bool)
	for rows.Next() {
		err = rows.Scan(&synonymName, &datasetName)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[synonymName] {
			for i, syn := range synonyms {
				if syn.Name == synonymName && !stringInSlice(datasetName, syn.Datasets) {
					synonyms[i].Datasets = append(synonyms[i].Datasets, datasetName)
					break
				}
			}
		} else {
			synonym.Name = synonymName
			var newSynDatasets []string
			synonym.Datasets = append(newSynDatasets, datasetName)
			synonyms = append(synonyms, synonym)
			exists[synonymName] = true
		}
	}
	drug.SYNS = synonyms

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data": drug,
			"type": "drug",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": drug,
			"type": "drug",
		})
	}
}

// DrugCells returns a list of cell lines tested against a drug, and
// number of experiments carried out with each cell line.
func DrugCells(c *gin.Context) {
	type DD struct {
		Cell     string   `json:"cell_line"`
		Datasets []string `json:"datasets"`
		Count    int      `json:"experiment-count"`
	}

	var (
		drugID      int
		cellName    string
		datasetName string
		experiment  DD
		experiments []DD
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT drug_id FROM drugs WHERE "
	var SQL2 string
	if searchByName(searchType) {
		SQL2 = "drug_name LIKE ?;"
	} else {
		SQL2 = "drug_id LIKE ?;"
	}
	SQL := SQL1 + SQL2
	row := db.QueryRow(SQL, id)
	err = row.Scan(&drugID)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(c, nil, http.StatusNotFound, "Drug not found in database")
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	q1 := "SELECT c.cell_name, da.dataset_name FROM experiments e JOIN cells c ON c.cell_id = e.cell_id "
	q2 := "JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.drug_id = ?;"
	query := q1 + q2
	rows, err := db.Query(query, drugID)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[string]bool)
	count := 0
	for rows.Next() {
		err = rows.Scan(&cellName, &datasetName)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[cellName] {
			for i, exp := range experiments {
				if exp.Cell == cellName {
					if !stringInSlice(datasetName, exp.Datasets) {
						experiments[i].Datasets = append(experiments[i].Datasets, datasetName)
					}
					experiments[i].Count++
					break
				}
			}
		} else {
			experiment.Cell = cellName
			var newExpDatasets []string
			experiment.Datasets = append(newExpDatasets, datasetName)
			experiment.Count = 1
			experiments = append(experiments, experiment)
			exists[cellName] = true
		}
		count++
	}

	if count == 0 {
		handleError(c, nil, http.StatusNotFound, "No cell lines found tested against this drug")
		return
	}

	desc := "List of cell lines tested against drug, and number of experiments carried out"

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        experiments,
			"description": desc,
			"total":       len(experiments),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        experiments,
			"description": desc,
			"total":       len(experiments),
		})
	}
}

// DrugTissues returns a list of tissues tested against a drug, and
// number of experiments carried out with each tissue.
func DrugTissues(c *gin.Context) {
	type DD struct {
		Tissue   string   `json:"tissue"`
		Datasets []string `json:"datasets"`
		Count    int      `json:"experiment-count"`
	}

	var (
		drugID      int
		tissueName  string
		datasetName string
		experiment  DD
		experiments []DD
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT drug_id FROM drugs WHERE "
	var SQL2 string
	if searchByName(searchType) {
		SQL2 = "drug_name LIKE ?;"
	} else {
		SQL2 = "drug_id LIKE ?;"
	}
	SQL := SQL1 + SQL2
	row := db.QueryRow(SQL, id)
	err = row.Scan(&drugID)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(c, nil, http.StatusNotFound, "Drug not found in database")
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	q1 := "SELECT t.tissue_name, da.dataset_name FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id "
	q2 := "JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.drug_id = ?;"
	query := q1 + q2
	rows, err := db.Query(query, drugID)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[string]bool)
	count := 0
	for rows.Next() {
		err = rows.Scan(&tissueName, &datasetName)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[tissueName] {
			for i, exp := range experiments {
				if exp.Tissue == tissueName {
					if !stringInSlice(datasetName, exp.Datasets) {
						experiments[i].Datasets = append(experiments[i].Datasets, datasetName)
					}
					experiments[i].Count++
					break
				}
			}
		} else {
			experiment.Tissue = tissueName
			var newExpDatasets []string
			experiment.Datasets = append(newExpDatasets, datasetName)
			experiment.Count = 1
			experiments = append(experiments, experiment)
			exists[tissueName] = true
		}
		count++
	}

	if count == 0 {
		handleError(c, nil, http.StatusNotFound, "No tissues found tested against this drug")
		return
	}

	desc := "List of tissues tested against drug, and number of experiments carried out"

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        experiments,
			"description": desc,
			"total":       len(experiments),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        experiments,
			"description": desc,
			"total":       len(experiments),
		})
	}
}
