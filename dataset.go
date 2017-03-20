package main

import (
	"net/http"

	"github.com/getsentry/raven-go"
	"gopkg.in/gin-gonic/gin.v1"
)

// DatasetReduced is a dataset with only two attributes
type DatasetReduced struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Dataset is a dataset datatype
type Dataset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetDatasets handles GET requests for /datasets
func GetDatasets(c *gin.Context) {
	var (
		dataset  DatasetReduced
		datasets []DatasetReduced
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

	rows, err := db.Query("select dataset_id, dataset_name from datasets;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&dataset.ID, &dataset.Name)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}
		datasets = append(datasets, dataset)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"count":    len(datasets),
		"datasets": datasets,
	})
}

// GetDatasetStats handles GET requests for /datasets/stats
func GetDatasetStats(c *gin.Context) {
	var (
		cstat DatasetStat
		tstat DatasetStat
		dstat DatasetStat
		estat DatasetStat

		cstats []DatasetStat
		tstats []DatasetStat
		dstats []DatasetStat
		estats []DatasetStat
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

	rows, err := db.Query("select dataset_id, cell_lines, tissues, drugs, experiments from dataset_statistics;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&cstat.Dataset, &cstat.Count, &tstat.Count, &dstat.Count, &estat.Count)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
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
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"statistics": gin.H{
			"cell_lines":  cstats,
			"tissues":     tstats,
			"drugs":       dstats,
			"experiments": estats,
		},
	})
}
