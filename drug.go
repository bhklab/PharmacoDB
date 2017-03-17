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

// GetDrugs handles GET requests for /drugs endpoint
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
		ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
		c.Abort()
		return
	}

	rows, err := db.Query("select drug_id, drug_name from drugs;")
	if err != nil {
		raven.CaptureError(err, nil)
		ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
		c.Abort()
		return
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			raven.CaptureError(err, nil)
			ErrorHandler(c, http.StatusInternalServerError, "Internal server error")
			c.Abort()
			return
		}
		drugs = append(drugs, drug)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"count": len(drugs),
		"data":  drugs,
	})
}
