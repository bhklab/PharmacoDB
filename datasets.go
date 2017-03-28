package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	null "gopkg.in/guregu/null.v3"
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

// getDatasetCells finds all cell lines tested in a dataset.
func getDatasetTypes(c *gin.Context, dtype string, iden string, queryStr string) {
	var (
		id        int
		name      null.String
		datatypes []DataTypeReduced
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	rows, err := db.Query(queryStr, iden)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[null.String]bool)
	iter := 0
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if !exists[name] {
			var datatype DataTypeReduced
			datatype.ID = id
			datatype.Name = name
			datatypes = append(datatypes, datatype)
			exists[name] = true
		}
		iter = 1
	}
	if iter == 0 {
		handleError(c, nil, http.StatusNotFound, fmt.Sprintf("No %s found tested in this cell line", dtype))
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": fmt.Sprintf("List of %s tested in dataset", dtype),
		"count":       len(datatypes),
		"data":        datatypes,
	})
}

// GetDatasetCellsByID handles GET requests for /datasets/ids/:id/cell_lines endpoint.
func GetDatasetCellsByID(c *gin.Context) {
	queryStr := "select c.cell_id, c.cell_name from datasets d inner join experiments e on e.dataset_id = d.dataset_id inner join cells c on c.cell_id = e.cell_id where d.dataset_id = ?"
	getDatasetTypes(c, "cell_lines", c.Param("id"), queryStr)
}

// GetDatasetCellsByName handles GET requests for /datasets/names/:name/cell_lines endpoint.
func GetDatasetCellsByName(c *gin.Context) {
	queryStr := "select c.cell_id, c.cell_name from datasets d inner join experiments e on e.dataset_id = d.dataset_id inner join cells c on c.cell_id = e.cell_id where d.dataset_name = ?"
	getDatasetTypes(c, "cell_lines", c.Param("name"), queryStr)
}

// GetDatasetTissuesByID handles GET requests for /datasets/ids/:id/tissues endpoint.
func GetDatasetTissuesByID(c *gin.Context) {
	queryStr := "select t.tissue_id, t.tissue_name from datasets d inner join experiments e on e.dataset_id = d.dataset_id inner join tissues t on t.tissue_id = e.tissue_id where d.dataset_id = ?"
	getDatasetTypes(c, "tissues", c.Param("id"), queryStr)
}

// GetDatasetTissuesByName handles GET requests for /datasets/names/:name/tissues endpoint.
func GetDatasetTissuesByName(c *gin.Context) {
	queryStr := "select t.tissue_id, t.tissue_name from datasets d inner join experiments e on e.dataset_id = d.dataset_id inner join tissues t on t.tissue_id = e.tissue_id where d.dataset_name = ?"
	getDatasetTypes(c, "tissues", c.Param("name"), queryStr)
}

// GetDatasetDrugsByID handles GET requests for /datasets/ids/:id/drugs endpoint.
func GetDatasetDrugsByID(c *gin.Context) {
	queryStr := "select dr.drug_id, dr.drug_name from datasets da inner join experiments e on e.dataset_id = da.dataset_id inner join drugs dr on dr.drug_id = e.drug_id where da.dataset_id = ?"
	getDatasetTypes(c, "drugs", c.Param("id"), queryStr)
}

// GetDatasetDrugsByName handles GET requests for /datasets/names/:name/drugs endpoint.
func GetDatasetDrugsByName(c *gin.Context) {
	queryStr := "select dr.drug_id, dr.drug_name from datasets da inner join experiments e on e.dataset_id = da.dataset_id inner join drugs dr on dr.drug_id = e.drug_id where da.dataset_name = ?"
	getDatasetTypes(c, "drugs", c.Param("name"), queryStr)
}
