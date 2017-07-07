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

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if isTrue, _ := strconv.ParseBool(c.DefaultQuery("all", "false")); isTrue {
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
		if shouldIndent {
			c.IndentedJSON(http.StatusOK, gin.H{
				"data":        tissues,
				"total":       len(tissues),
				"description": "List of all tissues in PharmacoDB",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":        tissues,
				"total":       len(tissues),
				"description": "List of all tissues in PharmacoDB",
			})
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	s := (page - 1) * limit
	selectSQL := "SELECT tissue_id, tissue_name FROM tissues"
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
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
	row := db.QueryRow("SELECT COUNT(*) FROM tissues;")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/tissues", page, total, limit)

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        tissues,
			"total":       total,
			"description": "List of all tissues in PharmacoDB",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        tissues,
			"total":       total,
			"description": "List of all tissues in PharmacoDB",
		})
	}
}
