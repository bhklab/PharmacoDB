package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Cell is a cell_line datatype.
type Cell struct {
	ID   int     `json:"id"`
	ACC  *string `json:"accession_id"`
	Name string  `json:"name"`
}

// IndexCell returns a list of all cell lines currently in database.
// Result is paginated by default, using: /cell_lines?page=int&per_page=int.
// To return all cell_lines in one call (without pagination), do: /cell_lines?all=true.
// Pagination links are available in response Link header.
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

	all := c.DefaultQuery("all", "false")
	if all == "true" {
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
		c.IndentedJSON(http.StatusOK, gin.H{
			"data":        cells,
			"total":       len(cells),
			"description": "List of all cell lines in PharmacoDB",
		})
		return
	}

	// Paginate response using page and per_page request values.
	// Default: page=1 and per_page=30.
	// Hence, /cell_lines is equivalent to /cell_lines?page=1&per_page=30 by default.

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
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS cell_id, accession_id, cell_name FROM cells limit %d,%d;", s, limit)
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
	row := db.QueryRow("SELECT FOUND_ROWS();")
	var total int
	err = row.Scan(&total)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Write pagination links in response header.

	var (
		prev    string
		prevRel string
		next    string
		nextRel string
	)
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	first := fmt.Sprintf("<https://api.pharmacodb.com/v1/cell_lines?page=%d&per_page=%d>", 1, limit)
	if (page > 1) && (page <= lastPage) {
		prev = fmt.Sprintf("<https://api.pharmacodb.com/v1/cell_lines?page=%d&per_page=%d>", page-1, limit)
		prevRel = "; rel=\"prev\", "
	}
	if (page >= 1) && (page < lastPage) {
		next = fmt.Sprintf("<https://api.pharmacodb.com/v1/cell_lines?page=%d&per_page=%d>", page+1, limit)
		nextRel = "; rel=\"next\", "
	}
	last := fmt.Sprintf("<https://api.pharmacodb.com/v1/cell_lines?page=%d&per_page=%d>", lastPage, limit)

	linknp := prev + prevRel + next + nextRel
	linkfl := first + "; rel=\"first\", " + last + "; rel=\"last\""
	link := linknp + linkfl

	c.Writer.Header().Set("Link", link)

	c.IndentedJSON(http.StatusOK, gin.H{
		"data":        cells,
		"page":        page,
		"per_page":    limit,
		"total":       total,
		"description": "List of all cell lines in PharmacoDB",
	})
}
