package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Drug is a drug datatype.
type Drug struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IndexDrug returns a list of all drugs currently in database.
func IndexDrug(c *gin.Context) {
	var (
		drug  Drug
		drugs []Drug
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	all := c.DefaultQuery("all", "false")
	if isTrue, _ := strconv.ParseBool(all); isTrue {
		rows, er := db.Query("SELECT drug_id, drug_name FROM drugs;")
		defer rows.Close()
		if er != nil {
			handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		for rows.Next() {
			err = rows.Scan(&drug.ID, &drug.Name)
			if err != nil {
				handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			drugs = append(drugs, drug)
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        drugs,
			"total":       len(drugs),
			"description": "List of all drugs in PharmacoDB",
		})
		return
	}

	curPage := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("per_page", "30")

	page, err := strconv.Atoi(curPage)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	limit, err := strconv.Atoi(perPage)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	s := (page - 1) * limit
	selectSQL := "SELECT drug_id, drug_name FROM drugs"
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		drugs = append(drugs, drug)
	}
	row := db.QueryRow("SELECT COUNT(*) FROM drugs;")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/drugs", page, total, limit)

	c.IndentedJSON(http.StatusOK, gin.H{
		"data":        drugs,
		"total":       total,
		"description": "List of all drugs in PharmacoDB",
	})
}
