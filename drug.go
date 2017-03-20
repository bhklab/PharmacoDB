package main

import (
	"net/http"

	"github.com/getsentry/raven-go"
	"gopkg.in/gin-gonic/gin.v1"
)

// DrugReduced is a drug with only two attributes
type DrugReduced struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Drug is a drug datatype
type Drug struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetDrugs handles GET requests for /drugs
func GetDrugs(c *gin.Context) {
	var (
		drug  DrugReduced
		drugs []DrugReduced
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

	rows, err := db.Query("select drug_id, drug_name from drugs;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}
		drugs = append(drugs, drug)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"count": len(drugs),
		"drugs": drugs,
	})
}

// GetDrugStats handles GET requests for /drugs/stats
func GetDrugStats(c *gin.Context) {
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

	rows, err := db.Query("select dataset_id, drugs from dataset_statistics;")
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
