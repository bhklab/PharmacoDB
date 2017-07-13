package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Cell is a cell_line datatype.
type Cell struct {
	ID     int       `json:"id"`
	ACC    *string   `json:"accession_id,omitempty"`
	Name   string    `json:"name"`
	Tissue *Tissue   `json:"tissue,omitempty"`
	SYNS   []Synonym `json:"synonyms,omitempty"`
}

// IndexCell returns a list of all cell lines currently in database.
func IndexCell(c *gin.Context) {
	var (
		cell  Cell
		cells []Cell
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if isTrue, _ := strconv.ParseBool(c.DefaultQuery("all", "false")); isTrue {
		rows, er := db.Query("SELECT cell_id, accession_id, cell_name FROM cells;")
		defer rows.Close()
		if er != nil {
			handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		for rows.Next() {
			err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
			if err != nil {
				handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			cells = append(cells, cell)
		}
		if shouldIndent {
			c.IndentedJSON(http.StatusOK, gin.H{
				"data":        cells,
				"total":       len(cells),
				"description": "List of all cell lines in PharmacoDB",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":        cells,
				"total":       len(cells),
				"description": "List of all cell lines in PharmacoDB",
			})
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	s := (page - 1) * limit
	selectSQL := "SELECT cell_id, accession_id, cell_name FROM cells"
	query := fmt.Sprintf("%s limit %d,%d;", selectSQL, s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		cells = append(cells, cell)
	}
	row := db.QueryRow("SELECT COUNT(*) FROM cells;")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.
	writeHeaderLinks(c, "/cell_lines", page, total, limit)

	if shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        cells,
			"total":       total,
			"description": "List of all cell lines in PharmacoDB",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        cells,
			"total":       total,
			"description": "List of all cell lines in PharmacoDB",
		})
	}
}

// ShowCell returns a cell line using ID, Name or Accession ID.
func ShowCell(c *gin.Context) {
	var (
		cell     Cell
		synonym  Synonym
		synonyms []Synonym
	)
	tissue := &Tissue{}

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name "
	SQL2 := "FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE "
	var SQL3 string
	if searchByName(searchType) {
		SQL3 = "c.cell_name LIKE ?;"
	} else if searchType == "accession" {
		SQL3 = "c.accession_id LIKE ?;"
	} else {
		SQL3 = "c.cell_id LIKE ?;"
	}
	SQL := SQL1 + SQL2 + SQL3
	row := db.QueryRow(SQL, id)
	err = row.Scan(&cell.ID, &cell.ACC, &cell.Name, &tissue.ID, &tissue.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(c, nil, http.StatusNotFound, "Cell line not found in database")
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}
	cell.Tissue = tissue

	q1 := "SELECT s.cell_name, d.dataset_name FROM source_cell_names s "
	q2 := "JOIN datasets d ON d.dataset_id = s.source_id WHERE s.cell_id = ?;"
	query := q1 + q2
	rows, err := db.Query(query, cell.ID)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	var (
		synonymName string
		datasetName string
	)
	exists := make(map[string]bool)
	for rows.Next() {
		err = rows.Scan(&synonymName, &datasetName)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[synonymName] {
			for i, syn := range synonyms {
				if syn.Name == synonymName && !stringInSlice(datasetName, syn.Datasets) {
					synonyms[i].Datasets = append(synonyms[i].Datasets, datasetName)
					break
				}
			}
		} else {
			synonym.Name = synonymName
			var newSynDatasets []string
			synonym.Datasets = append(newSynDatasets, datasetName)
			synonyms = append(synonyms, synonym)
			exists[synonymName] = true
		}
	}
	cell.SYNS = synonyms

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data": cell,
			"type": "cell_line",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": cell,
			"type": "cell_line",
		})
	}
}

// CellDrugs returns a list of drugs tested with a cell line, and
// number of experiments carried out with each drug.
func CellDrugs(c *gin.Context) {
	var (
		cellID      int
		drugName    string
		datasetName string
		experiment  DD
		experiments []DD
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	id := c.Param("id")
	searchType := c.DefaultQuery("type", "id")

	SQL1 := "SELECT cell_id FROM cells WHERE "
	var SQL2 string
	if searchByName(searchType) {
		SQL2 = "cell_name LIKE ?;"
	} else if searchType == "accession" {
		SQL2 = "accession_id LIKE ?;"
	} else {
		SQL2 = "cell_id LIKE ?;"
	}
	SQL := SQL1 + SQL2
	row := db.QueryRow(SQL, id)
	err = row.Scan(&cellID)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(c, nil, http.StatusNotFound, "Cell line not found in database")
		} else {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	q1 := "SELECT d.drug_name, da.dataset_name FROM experiments e JOIN drugs d ON d.drug_id = e.drug_id "
	q2 := "JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.cell_id = ?"
	query := q1 + q2
	rows, err := db.Query(query, cellID)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[string]bool)
	count := 0
	for rows.Next() {
		err = rows.Scan(&drugName, &datasetName)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[drugName] {
			for i, exp := range experiments {
				if exp.Drug == drugName && !stringInSlice(datasetName, exp.Datasets) {
					experiments[i].Datasets = append(experiments[i].Datasets, datasetName)
					experiments[i].Count++
					break
				}
			}
		} else {
			experiment.Drug = drugName
			var newExpDatasets []string
			experiment.Datasets = append(newExpDatasets, datasetName)
			experiment.Count = 1
			experiments = append(experiments, experiment)
			exists[drugName] = true
		}
		count++
	}

	if count == 0 {
		handleError(c, nil, http.StatusNotFound, "No drugs found tested with this cell line")
		return
	}

	desc := "List of drugs tested with cell line, and number of experiments carried out with each drug"

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        experiments,
			"description": desc,
			"total":       len(experiments),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":        experiments,
			"description": desc,
			"total":       len(experiments),
		})
	}
}
