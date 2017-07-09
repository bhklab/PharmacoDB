package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Dataset is a dataset datatype.
type Dataset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IndexDataset returns a list of all datasets currently in database.
func IndexDataset(c *gin.Context) {
	var (
		dataset  Dataset
		datasets []Dataset
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if isTrue, _ := strconv.ParseBool(c.DefaultQuery("all", "false")); isTrue {
		rows, er := db.Query("SELECT dataset_id, dataset_name FROM datasets;")
		defer rows.Close()
		if er != nil {
			handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		for rows.Next() {
			err = rows.Scan(&dataset.ID, &dataset.Name)
			if err != nil {
				handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			datasets = append(datasets, dataset)
		}
		if shouldIndent {
			c.IndentedJSON(http.StatusOK, gin.H{
				"data":        datasets,
				"total":       len(datasets),
				"description": "List of all datasets in PharmacoDB",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":        datasets,
				"total":       len(datasets),
				"description": "List of all datasets in PharmacoDB",
			})
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	s := (page - 1) * limit
	selectSQL := "SELECT dataset_id, dataset_name FROM datasets"
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&dataset.ID, &dataset.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		datasets = append(datasets, dataset)
	}
	row := db.QueryRow("SELECT COUNT(*) FROM datasets;")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/datasets", page, total, limit)

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        datasets,
			"total":       total,
			"description": "List of all datasets in PharmacoDB",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        datasets,
			"total":       total,
			"description": "List of all datasets in PharmacoDB",
		})
	}
}

// ShowDataset returns a dataset using ID or Name.
func ShowDataset(c *gin.Context) {
	var dataset Dataset

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT dataset_id, dataset_name FROM datasets WHERE "
	var SQL2 string
	if searchByName(searchType) {
		SQL2 = "dataset_name LIKE ?;"
	} else {
		SQL2 = "dataset_id LIKE ?;"
	}
	SQL := SQL1 + SQL2
	row := db.QueryRow(SQL, id)
	err = row.Scan(&dataset.ID, &dataset.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			message := fmt.Sprintf("Dataset with ID:%s not found in database", id)
			handleError(c, nil, http.StatusNotFound, message)
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data": dataset,
			"type": "dataset",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": dataset,
			"type": "dataset",
		})
	}
}
