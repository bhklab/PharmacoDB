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

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if isTrue, _ := strconv.ParseBool(c.DefaultQuery("all", "false")); isTrue {
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
		if shouldIndent {
			c.IndentedJSON(http.StatusOK, gin.H{
				"data":        drugs,
				"total":       len(drugs),
				"description": "List of all drugs in PharmacoDB",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":        drugs,
				"total":       len(drugs),
				"description": "List of all drugs in PharmacoDB",
			})
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

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

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        drugs,
			"total":       total,
			"description": "List of all drugs in PharmacoDB",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        drugs,
			"total":       total,
			"description": "List of all drugs in PharmacoDB",
		})
	}
}
