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

// GetDatasets handles GET requests for /datasets endpoint
func GetDatasets(c *gin.Context) {
	var (
		dataset  DatasetReduced
		datasets []DatasetReduced
	)

	db := InitDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
		c.Abort()
		return
	}

	rows, err := db.Query("select dataset_id, dataset_name from datasets;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&dataset.ID, &dataset.Name)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
			c.Abort()
			return
		}
		datasets = append(datasets, dataset)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"count": len(datasets),
		"data":  datasets,
	})
}
