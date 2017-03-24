package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetDatasets handles GET requests for /datasets endpoint.
func GetDatasets(c *gin.Context) {
	getDataTypes(c, "List of all datasets in pharmacodb", "select dataset_id, dataset_name from datasets;")
}

// GetDatasetStats handles GET requests for /datasets/stats
func GetDatasetStats(c *gin.Context) {
	var (
		cstat  DatasetStat
		tstat  DatasetStat
		dstat  DatasetStat
		estat  DatasetStat
		cstats []DatasetStat
		tstats []DatasetStat
		dstats []DatasetStat
		estats []DatasetStat
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	rows, err := db.Query("select dataset_id, cell_lines, tissues, drugs, experiments from dataset_statistics;")
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&cstat.Dataset, &cstat.Count, &tstat.Count, &dstat.Count, &estat.Count)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		tstat.Dataset = cstat.Dataset
		dstat.Dataset = cstat.Dataset
		estat.Dataset = cstat.Dataset
		cstats = append(cstats, cstat)
		tstats = append(tstats, tstat)
		dstats = append(dstats, dstat)
		estats = append(estats, estat)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": "Number of items tested in each dataset per datatype, as well as number of experiments carried out in each dataset",
		"data": gin.H{
			"cell_lines":  cstats,
			"tissues":     tstats,
			"drugs":       dstats,
			"experiments": estats,
		},
	})
}

// GetDatasetIDs handles GET requests for /datasets/ids endpoint.
func GetDatasetIDs(c *gin.Context) {
	getDataTypeIDs(c, "List of all dataset IDs in pharmacodb", "select dataset_id from datasets;")
}

// GetDatasetNames handles GET requests for /datasets/names endpoint.
func GetDatasetNames(c *gin.Context) {
	getDataTypeNames(c, "List of all dataset names in pharmacodb", "select dataset_name from datasets;")
}

// getDataset finds a dataset with either ID or name.
func getDataset(c *gin.Context, ptype string) {
	var (
		dataset  Dataset
		queryStr string
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	iden := c.Param(ptype)
	if ptype == "id" {
		queryStr = "select dataset_id, dataset_name from datasets where dataset_id = ?"
	} else {
		queryStr = "select dataset_id, dataset_name from datasets where dataset_name = ?"
	}
	row := db.QueryRow(queryStr, iden)
	err = row.Scan(&dataset.ID, &dataset.Name)
	if err == sql.ErrNoRows {
		handleError(c, nil, http.StatusNotFound, fmt.Sprintf("Dataset with %s - %s - not found in pharmacodb", ptype, iden))
		return
	} else if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"type": "dataset",
		"data": dataset,
	})
}

// GetDatasetByID handles GET requests for /datasets/ids/:id endpoints.
func GetDatasetByID(c *gin.Context) {
	getDataset(c, "id")
}

// GetDatasetByName handles GET requests for /datasets/names/:name endpoints.
func GetDatasetByName(c *gin.Context) {
	getDataset(c, "name")
}
