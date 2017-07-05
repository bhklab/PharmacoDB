package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Tissue is a tissue datatype.
type Tissue struct {
	ID   int     `json:"id"`
	Name *string `json:"name,omitempty"`
}

// IndexTissue returns a list of all tissues currently in database.
func IndexTissue(c *gin.Context) {
	var (
		tissue  Tissue
		tissues []Tissue
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	all := c.DefaultQuery("all", "false")
	if isTrue, _ := strconv.ParseBool(all); isTrue {
		rows, er := db.Query("SELECT tissue_id, tissue_name FROM tissues;")
		defer rows.Close()
		if er != nil {
			handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		for rows.Next() {
			err = rows.Scan(&tissue.ID, &tissue.Name)
			if err != nil {
				handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			tissues = append(tissues, tissue)
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        tissues,
			"total":       len(tissues),
			"description": "List of all tissues in PharmacoDB",
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
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS tissue_id, tissue_name FROM tissues limit %d,%d;", s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		tissues = append(tissues, tissue)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/tissues", page, total, limit)

	c.IndentedJSON(http.StatusOK, gin.H{
		"data":        tissues,
		"page":        page,
		"per_page":    limit,
		"total":       total,
		"description": "List of all tissues in PharmacoDB",
	})
}
