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
