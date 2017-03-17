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

// GetTissues handles GET requests for /tissues endpoint
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
		ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
		c.Abort()
		return
	}

	rows, err := db.Query("select tissue_id, tissue_name from tissues;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
			c.Abort()
			return
		}
		tissues = append(tissues, tissue)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"count": len(tissues),
		"data":  tissues,
	})
}
