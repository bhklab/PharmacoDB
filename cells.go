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

// IndexCell returns a list of all cell lines available in database (paginated by default).
// The request responds to a url matching: /cell_lines?all=&page=&per_page=
// To return all cell_lines in one call, do: /cell_lines?all=true
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

	// Paginate response using page and per_page values.
	// Default: page=1 and per_page=30
	// If user does not explicitely give values, default returned.
	// Hence, /cell_lines response is the same as /cell_lines?page=1&per_page=30

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

	// Define pagination links using page and limit.
	var (
		prev string
		next string
	)
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	first := fmt.Sprintf("<https://api.pharmacodb.com/v1/cell_lines?page=1&per_page=%d>; rel=\"first\", ", limit)
	if page != 1 {
		prev = fmt.Sprintf("<https://api.pharmacodb.com/v1/cell_lines?page=%d&per_page=%d>; rel=\"prev\", ", page-1, limit)
	}
	if page != lastPage {
		next = fmt.Sprintf("<https://api.pharmacodb.com/v1/cell_lines?page=%d&per_page=%d>; rel=\"next\", ", page+1, limit)
	}

	last := fmt.Sprintf("<https://api.pharmacodb.com/v1/cell_lines?page=%d&per_page=%d>; rel=\"last\"", lastPage, limit)
	link := first + prev + next + last
	// Write links to response header.
	c.Writer.Header().Set("Link", link)

	c.IndentedJSON(http.StatusOK, gin.H{
		"data":        cells,
		"page":        page,
		"per_page":    limit,
		"total":       total,
		"description": "List of all cell lines in PharmacoDB",
	})
}
