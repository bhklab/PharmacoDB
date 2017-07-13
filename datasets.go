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
				"data":  datasets,
				"total": len(datasets),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":  datasets,
				"total": len(datasets),
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
			"data":  datasets,
			"total": total,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":  datasets,
			"total": total,
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

// DatasetCells returns a list of all cell lines tested in a particular dataset.
func DatasetCells(c *gin.Context) {
	var (
		datasetID   int
		datasetName string
		cell        Cell
		cells       []Cell
	)

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
	err = row.Scan(&datasetID, &datasetName)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(c, nil, http.StatusNotFound, "Dataset not found in database")
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if isTrue, _ := strconv.ParseBool(c.DefaultQuery("all", "false")); isTrue {
		q1 := "SELECT c.cell_id, c.cell_name FROM experiments e "
		q2 := "JOIN cells c ON c.cell_id = e.cell_id WHERE e.dataset_id = ? GROUP BY e.cell_id;"
		query := q1 + q2
		rows, er := db.Query(query, datasetID)
		defer rows.Close()
		if er != nil {
			handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		count := 0
		for rows.Next() {
			err = rows.Scan(&cell.ID, &cell.Name)
			if err != nil {
				handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			cells = append(cells, cell)
			count++
		}
		if count == 0 {
			handleError(c, nil, http.StatusNotFound, "No cell lines found tested in this dataset")
			return
		}
		if shouldIndent {
			c.IndentedJSON(http.StatusOK, gin.H{
				"data":  cells,
				"total": len(cells),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":  cells,
				"total": len(cells),
			})
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	s := (page - 1) * limit
	q1 := "SELECT c.cell_id, c.cell_name FROM experiments e "
	q2 := "JOIN cells c ON c.cell_id = e.cell_id WHERE e.dataset_id = ? GROUP BY e.cell_id"
	selectSQL := q1 + q2
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
	rows, err := db.Query(query, id)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		cells = append(cells, cell)
	}
	row = db.QueryRow("SELECT COUNT(DISTINCT cell_id) FROM experiments WHERE dataset_id = ?;", datasetID)
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/datasets/:id/cell_lines", page, total, limit)

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":  cells,
			"total": total,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":  cells,
			"total": total,
		})
	}
}

// DatasetDrugs returns a list of all drugs tested in a dataset.
func DatasetDrugs(c *gin.Context) {
	var (
		datasetID int
		drug      Drug
		drugs     []Drug
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT dataset_id FROM datasets WHERE "
	var SQL2 string
	if searchByName(searchType) {
		SQL2 = "dataset_name LIKE ?;"
	} else {
		SQL2 = "dataset_id LIKE ?;"
	}
	SQL := SQL1 + SQL2
	row := db.QueryRow(SQL, id)
	err = row.Scan(&datasetID)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(c, nil, http.StatusNotFound, "Dataset not found in database")
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if isTrue, _ := strconv.ParseBool(c.DefaultQuery("all", "false")); isTrue {
		q1 := "SELECT d.drug_id, d.drug_name FROM experiments e "
		q2 := "JOIN drugs d ON d.drug_id = e.drug_id WHERE e.dataset_id = ? GROUP BY e.drug_id;"
		query := q1 + q2
		rows, er := db.Query(query, datasetID)
		defer rows.Close()
		if er != nil {
			handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		count := 0
		for rows.Next() {
			err = rows.Scan(&drug.ID, &drug.Name)
			if err != nil {
				handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			drugs = append(drugs, drug)
			count++
		}
		if count == 0 {
			handleError(c, nil, http.StatusNotFound, "No drugs found tested in this dataset")
			return
		}
		if shouldIndent {
			c.IndentedJSON(http.StatusOK, gin.H{
				"data":  drugs,
				"total": len(drugs),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":  drugs,
				"total": len(drugs),
			})
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	s := (page - 1) * limit
	q1 := "SELECT d.drug_id, d.drug_name FROM experiments e "
	q2 := "JOIN drugs d ON d.drug_id = e.drug_id WHERE e.dataset_id = ? GROUP BY e.drug_id"
	selectSQL := q1 + q2
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
	rows, err := db.Query(query, id)
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
	row = db.QueryRow("SELECT COUNT(DISTINCT drug_id) FROM experiments WHERE dataset_id = ?;", datasetID)
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/datasets/:id/cell_lines", page, total, limit)

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":  drugs,
			"total": total,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":  drugs,
			"total": total,
		})
	}
}
