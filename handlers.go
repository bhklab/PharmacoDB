package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Count returns the total number of records in table,
// or error in case of failure.
func Count(table string) (int, error) {
	var count int
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return 0, err
	}
	query := "SELECT COUNT(*) FROM " + table + ";"
	row := db.QueryRow(query)
	err = row.Scan(&count)
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return 0, err
	}
	return count, nil
}

// RenderJSON outputs response as either indented or non-indented
// depending on setting by parent function.
func RenderJSON(c *gin.Context, indent bool, obj interface{}) {
	if indent {
		c.IndentedJSON(http.StatusOK, obj)
	} else {
		c.JSON(http.StatusOK, obj)
	}
}

// writeHeaderLinks writes pagination links in response header.
// Links available under 'Link' header, including (first, prev, next, last).
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

// CellsHandler is a handler for '/cell_lines' endpoint.
// Lists all cell lines in database (paginated and non-paginated).
func CellsHandler(c *gin.Context) {
	// Optional parameters
	// Page and limit are used for paginated response (default).
	// If listAll is set to true, it takes precedence over page and limit,
	//    returning all cell lines (without pagination).
	// Response indented by default, can be set to false for non-indented responses.
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	if listAll {
		cells, err := ListAllCells()
		if err != nil {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		RenderJSON(c, indent, gin.H{"data": cells, "total": len(cells)})
		return
	}
	cells, err := ListPaginatedCells(page, limit)
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
	RenderJSON(c, indent, gin.H{"data": cells, "total": count})
}

// TissuesHandler is a handler for '/tissues' endpoint.
// Lists all tissues in database (paginated and non-paginated).
func TissuesHandler(c *gin.Context) {
	// Optional parameters (see CellsHandler)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	if listAll {
		tissues, err := ListAllTissues()
		if err != nil {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		RenderJSON(c, indent, gin.H{"data": tissues, "total": len(tissues)})
		return
	}
	tissues, err := ListPaginatedTissues(page, limit)
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	count, err := Count("tissues")
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// Write pagination links in response header.
	writeHeaderLinks(c, "/tissues", page, count, limit)
	RenderJSON(c, indent, gin.H{"data": tissues, "total": count})
}

// DrugsHandler is a handler for '/drugs' endpoint.
// Lists all drugs in database (paginated and non-paginated).
func DrugsHandler(c *gin.Context) {
	// Optional parameters (see CellsHandler)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	if listAll {
		drugs, err := ListAllDrugs()
		if err != nil {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		RenderJSON(c, indent, gin.H{"data": drugs, "total": len(drugs)})
		return
	}
	drugs, err := ListPaginatedDrugs(page, limit)
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	count, err := Count("drugs")
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// Write pagination links in response header.
	writeHeaderLinks(c, "/drugs", page, count, limit)
	RenderJSON(c, indent, gin.H{"data": drugs, "total": count})
}

// DatasetsHandler is a handler for '/datasets' endpoint.
// Lists all datasets in database (paginated and non-paginated).
func DatasetsHandler(c *gin.Context) {
	// Optional parameters (see CellsHandler)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	if listAll {
		datasets, err := ListAllDatasets()
		if err != nil {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		RenderJSON(c, indent, gin.H{"data": datasets, "total": len(datasets)})
		return
	}
	datasets, err := ListPaginatedDatasets(page, limit)
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	count, err := Count("datasets")
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// Write pagination links in response header.
	writeHeaderLinks(c, "/datasets", page, count, limit)
	RenderJSON(c, indent, gin.H{"data": datasets, "total": count})
}

// ExperimentsHandler is a handler for '/experiments' endpoint.
// Lists all experiments in database (paginated only).
func ExperimentsHandler(c *gin.Context) {
	// Optional parameters (see CellsHandler)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	// Set max limit per_page to 1000
	if limit > 1000 {
		limit = 1000
	}
	experiments, err := ListPaginatedExperiments(page, limit)
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	count, err := Count("experiments")
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// Write pagination links in response header.
	writeHeaderLinks(c, "/experiments", page, count, limit)
	RenderJSON(c, indent, gin.H{"data": experiments, "total": count})
}
