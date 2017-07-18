package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HandlerFunc is a gin HandlerFunc.
type HandlerFunc gin.HandlerFunc

// Count returns the total number of records in table,
// or error in case of failure.
func Count(table string) (int, error) {
	var count int
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return 0, err
	}
	query := "SELECT COUNT(*) FROM " + table
	row := db.QueryRow(query)
	err = row.Scan(&count)
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return 0, err
	}
	return count, nil
}

// CustomJSON outputs response as either indented or non-indented
// depending on setting by parent function.
func CustomJSON(c *gin.Context, obj gin.H, indent bool) {
	if indent {
		c.IndentedJSON(http.StatusOK, obj)
	} else {
		c.JSON(http.StatusOK, obj)
	}
}

// writeHeaderLinks writes pagination links in response header.
// Links available under 'Link' header, including (prev, next, first, last).
func writeHeaderLinks(c *gin.Context, endpoint string, page int, total int, limit int) {
	var (
		prev    string
		relPrev string
		next    string
		relNext string
	)
	pattern := "<https://api.pharmacodb.com/" + APIVersion() + "%s?page=%d&per_page=%d>"
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	first := fmt.Sprintf(pattern, endpoint, 1, limit)
	relFirst := "; rel=\"first\", "
	last := fmt.Sprintf(pattern, endpoint, lastPage, limit)
	relLast := "; rel=\"last\""
	if (page > 1) && (page <= lastPage) {
		prev = fmt.Sprintf(pattern, endpoint, page-1, limit)
		relPrev = "; rel=\"prev\", "
	}
	if (page >= 1) && (page < lastPage) {
		next = fmt.Sprintf(pattern, endpoint, page+1, limit)
		relNext = "; rel=\"next\", "
	}
	link := first + relFirst + prev + relPrev + next + relNext + last + relLast
	c.Writer.Header().Set("Link", link)
}

// CellsHandler is a handler for '/cell_lines'.
// Lists all cell lines in database.
func CellsHandler(c *gin.Context) {
	// Optional parameters
	// page and limit are used for paginated response (default).
	// If listAll is set to true, it takes precedence over page and limit,
	//    returning all cell lines (without pagination).
	// Response indented by default, can be set to false for non-indented responses.
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))

	if listAll {
		cells, err := NonPaginatedCells()
		if err != nil {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		CustomJSON(c, gin.H{"data": cells, "total": len(cells)}, indent)
		return
	}
	cells, err := PaginatedCells(page, limit)
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	count, err := Count("cells")
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// Write pagination links in response header.
	writeHeaderLinks(c, "/cell_lines", page, count, limit)
	CustomJSON(c, gin.H{"data": cells, "total": count}, indent)
}
