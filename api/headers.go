package api

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// WriteHeader writes metadata to response header.
func WriteHeader(c *gin.Context, endpoint string, page int, limit int, total int) {
	var (
		prev    string
		relPrev string
		next    string
		relNext string
	)

	pattern := "<https://api.pharmacodb.com/" + Version() + "%s?page=%d&per_page=%d>"
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

	// Write all custom headers.
	c.Writer.Header().Set("Link", link)
	c.Writer.Header().Set("Pagination-Current-Page", strconv.Itoa(page))
	c.Writer.Header().Set("Pagination-Last-Page", strconv.Itoa(lastPage))
	c.Writer.Header().Set("Pagination-Per-Page", strconv.Itoa(limit))
	c.Writer.Header().Set("Total-Records", strconv.Itoa(total))
}

// IndexCellsHEAD returns header info for cell lines.
// Handles HEAD requests for /cell_lines.
func IndexCellsHEAD(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	all, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))

	total, err := Count("cells")
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	if all {
		c.Writer.Header().Set("Total-Records", strconv.Itoa(total))
	} else {
		WriteHeader(c, "/cell_lines", page, limit, total)
	}
}
