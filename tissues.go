package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Tissue is a tissue datatype.
type Tissue struct {
	ID   int       `json:"id"`
	Name *string   `json:"name,omitempty"`
	SYNS []Synonym `json:"synonyms,omitempty"`
}

// IndexTissue returns a list of all tissues currently in database.
func IndexTissue(c *gin.Context) {
	var (
		tissue  Tissue
		tissues []Tissue
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if isTrue, _ := strconv.ParseBool(c.DefaultQuery("all", "false")); isTrue {
		rows, er := db.Query("SELECT tissue_id, tissue_name FROM tissues;")
		defer rows.Close()
		if er != nil {
			handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		for rows.Next() {
			err = rows.Scan(&tissue.ID, &tissue.Name)
			if err != nil {
				handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			tissues = append(tissues, tissue)
		}
		if shouldIndent {
			c.IndentedJSON(http.StatusOK, gin.H{
				"data":        tissues,
				"total":       len(tissues),
				"description": "List of all tissues in PharmacoDB",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":        tissues,
				"total":       len(tissues),
				"description": "List of all tissues in PharmacoDB",
			})
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	s := (page - 1) * limit
	selectSQL := "SELECT tissue_id, tissue_name FROM tissues"
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		tissues = append(tissues, tissue)
	}
	row := db.QueryRow("SELECT COUNT(*) FROM tissues;")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/tissues", page, total, limit)

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        tissues,
			"total":       total,
			"description": "List of all tissues in PharmacoDB",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        tissues,
			"total":       total,
			"description": "List of all tissues in PharmacoDB",
		})
	}
}

// ShowTissue returns a tissue using ID or Name.
func ShowTissue(c *gin.Context) {
	var (
		tissue   Tissue
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

	SQL1 := "SELECT tissue_id, tissue_name FROM tissues WHERE "
	var SQL2 string
	if searchByName(searchType) {
		SQL2 = "tissue_name LIKE ?;"
	} else {
		SQL2 = "tissue_id LIKE ?;"
	}
	SQL := SQL1 + SQL2
	row := db.QueryRow(SQL, id)
	err = row.Scan(&tissue.ID, &tissue.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			message := fmt.Sprintf("Tissue with ID:%s not found in database", id)
			handleError(c, nil, http.StatusNotFound, message)
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	q1 := "SELECT s.tissue_name, d.dataset_name FROM source_tissue_names s "
	q2 := "JOIN datasets d ON d.dataset_id = s.source_id WHERE s.tissue_id = ?;"
	query := q1 + q2
	rows, err := db.Query(query, tissue.ID)
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
	tissue.SYNS = synonyms

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data": tissue,
			"type": "tissue",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": tissue,
			"type": "tissue",
		})
	}
}

// TissueCells returns a list of all cell lines of a specific tissue type.
func TissueCells(c *gin.Context) {
	var (
		tissueID   int
		tissueName string
		cell       Cell
		cells      []Cell
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT tissue_id, tissue_name FROM tissues WHERE "
	var SQL2 string
	if searchByName(searchType) {
		SQL2 = "tissue_name LIKE ?;"
	} else {
		SQL2 = "tissue_id LIKE ?;"
	}
	SQL := SQL1 + SQL2
	row := db.QueryRow(SQL, id)
	err = row.Scan(&tissueID, &tissueName)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(c, nil, http.StatusNotFound, "Tissue not found in database")
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	query := "SELECT cell_id, cell_name FROM cells WHERE tissue_id = ?;"
	rows, err := db.Query(query, tissueID)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		cells = append(cells, cell)
	}

	desc := fmt.Sprintf("List of all cell lines of %s tissue type in PharmacoDB", tissueName)

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        cells,
			"total":       len(cells),
			"description": desc,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        cells,
			"total":       len(cells),
			"description": desc,
		})
	}
}

// TissueDrugs returns a list of drugs tested with a tissue, and
// number of experiments carried out with each drug.
func TissueDrugs(c *gin.Context) {
	type DD struct {
		Drug     string   `json:"drug"`
		Datasets []string `json:"datasets"`
		Count    int      `json:"experiment-count"`
	}

	var (
		tissueID    int
		drugName    string
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

	SQL1 := "SELECT tissue_id FROM tissues WHERE "
	var SQL2 string
	if searchByName(searchType) {
		SQL2 = "tissue_name LIKE ?;"
	} else {
		SQL2 = "tissue_id LIKE ?;"
	}
	SQL := SQL1 + SQL2
	row := db.QueryRow(SQL, id)
	err = row.Scan(&tissueID)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(c, nil, http.StatusNotFound, "Tissue not found in database")
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	q1 := "SELECT d.drug_name, da.dataset_name FROM experiments e JOIN drugs d ON d.drug_id = e.drug_id "
	q2 := "JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.tissue_id = ?"
	query := q1 + q2
	rows, err := db.Query(query, tissueID)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[string]bool)
	count := 0
	for rows.Next() {
		err = rows.Scan(&drugName, &datasetName)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[drugName] {
			for i, exp := range experiments {
				if exp.Drug == drugName {
					if !stringInSlice(datasetName, exp.Datasets) {
						experiments[i].Datasets = append(experiments[i].Datasets, datasetName)
					}
					experiments[i].Count++
					break
				}
			}
		} else {
			experiment.Drug = drugName
			var newExpDatasets []string
			experiment.Datasets = append(newExpDatasets, datasetName)
			experiment.Count = 1
			experiments = append(experiments, experiment)
			exists[drugName] = true
		}
		count++
	}

	if count == 0 {
		handleError(c, nil, http.StatusNotFound, "No drugs found tested with this tissue")
		return
	}

	desc := "List of drugs tested with tissue, and number of experiments carried out with each drug"

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
