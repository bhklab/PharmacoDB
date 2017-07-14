package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CellCountPerTissue returns the cell line count per tissue for all tissues.
func CellCountPerTissue(c *gin.Context) {
	type TissueCellCount struct {
		Tissue Tissue `json:"tissue"`
		Count  int    `json:"cell-line-count"`
	}
	var (
		tcc  TissueCellCount
		tccs []TissueCellCount
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	q1 := "SELECT t.tissue_id, t.tissue_name, (SELECT COUNT(DISTINCT c.cell_id) "
	q2 := "FROM cells c WHERE c.tissue_id = t.tissue_id) as cell_count FROM tissues t;"
	query := q1 + q2
	rows, er := db.Query(query)
	defer rows.Close()
	if er != nil {
		handleError(c, er, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for rows.Next() {
		err = rows.Scan(&tcc.Tissue.ID, &tcc.Tissue.Name, &tcc.Count)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		tccs = append(tccs, tcc)
	}

	if shouldIndent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true")); shouldIndent {
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":          tccs,
			"total-tissues": len(tccs),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":          tccs,
			"total-tissues": len(tccs),
		})
	}
}
