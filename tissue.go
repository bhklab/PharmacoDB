package main

import (
	"net/http"

	raven "github.com/getsentry/raven-go"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/guregu/null.v3"
)

// TissueReduced is a tissue with only two attributes
type TissueReduced struct {
	ID   int         `json:"id"`
	Name null.String `json:"name"`
}

// Tissue is a tissue datatype
type Tissue struct {
	ID   int         `json:"id"`
	Name null.String `json:"name"`
}

// GetTissues handles GET requests for /tissues
func GetTissues(c *gin.Context) {
	var (
		tissue  TissueReduced
		tissues []TissueReduced
	)

	db := InitDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}

	rows, err := db.Query("select tissue_id, tissue_name from tissues;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}
		tissues = append(tissues, tissue)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"count":   len(tissues),
		"tissues": tissues,
	})
}

// GetTissueStats handles GET requests for /tissues/stats
func GetTissueStats(c *gin.Context) {
	var (
		stat  DatasetStat
		stats []DatasetStat
	)

	db := InitDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}

	rows, err := db.Query("select dataset_id, tissues from dataset_statistics;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&stat.Dataset, &stat.Count)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}
		stats = append(stats, stat)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"statistics": stats,
	})
}
