package main

import (
	"database/sql"
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

// RenderJSONwithMeta outputs response as either indented or non-indented
// along with metadata about the reponse.
// Metadata includes current page, last page, per_page count and total count of result records.
func RenderJSONwithMeta(c *gin.Context, indent bool, page int, limit int, count int, include string, obj interface{}) {
	var data interface{}
	if include == "metadata" {
		lastPage := int(math.Ceil(float64(count) / float64(limit)))
		meta := gin.H{"page": page, "per_page": limit, "last_page": lastPage, "total": count}
		data = gin.H{"metadata": meta, "data": obj}
	} else {
		data = obj
	}
	if indent {
		c.IndentedJSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusOK, data)
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
	// Write all custom headers.
	c.Writer.Header().Set("Link", link)
	c.Writer.Header().Set("Pagination-Current-Page", strconv.Itoa(page))
	c.Writer.Header().Set("Pagination-Last-Page", strconv.Itoa(lastPage))
	c.Writer.Header().Set("Pagination-Per-Page", strconv.Itoa(limit))
	c.Writer.Header().Set("Total-Records", strconv.Itoa(total))
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
	include := c.Query("include") // shortcut for c.Request.URL.Query().Get("include")
	if listAll {
		cells, err := NonPaginatedCells()
		if err != nil {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(cells)))
		RenderJSON(c, indent, cells)
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
	RenderJSONwithMeta(c, indent, page, limit, count, include, cells)
}

// CellHandler is a handler for '/cell_lines/:id' endpoint.
// Returns a single cell line.
func CellHandler(c *gin.Context) {
	// Optional parameters
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	typ := c.DefaultQuery("type", "id")
	// path parameter
	id := c.Param("id")
	cell, err := FindCell(id, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			LogPublicError(c, ErrorTypePublic, http.StatusNotFound, "Cell Line Not Found")
		} else {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}
	err = cell.Annotate()
	if err != nil {
		LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
	}
	RenderJSON(c, indent, cell)
}

// TissuesHandler is a handler for '/tissues' endpoint.
// Lists all tissues in database (paginated and non-paginated).
func TissuesHandler(c *gin.Context) {
	// Optional parameters (see CellsHandler)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	if listAll {
		tissues, err := NonPaginatedTissues()
		if err != nil {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(tissues)))
		RenderJSON(c, indent, tissues)
		return
	}
	tissues, err := PaginatedTissues(page, limit)
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
	RenderJSONwithMeta(c, indent, page, limit, count, include, tissues)
}

// DrugsHandler is a handler for '/drugs' endpoint.
// Lists all drugs in database (paginated and non-paginated).
func DrugsHandler(c *gin.Context) {
	// Optional parameters (see CellsHandler)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	if listAll {
		drugs, err := NonPaginatedDrugs()
		if err != nil {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(drugs)))
		RenderJSON(c, indent, drugs)
		return
	}
	drugs, err := PaginatedDrugs(page, limit)
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
	RenderJSONwithMeta(c, indent, page, limit, count, include, drugs)
}

// DatasetsHandler is a handler for '/datasets' endpoint.
// Lists all datasets in database (paginated and non-paginated).
func DatasetsHandler(c *gin.Context) {
	// Optional parameters (see CellsHandler)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	if listAll {
		datasets, err := NonPaginatedDatasets()
		if err != nil {
			LogPublicError(c, ErrorTypePublic, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(datasets)))
		RenderJSON(c, indent, datasets)
		return
	}
	datasets, err := PaginatedDatasets(page, limit)
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
	RenderJSONwithMeta(c, indent, page, limit, count, include, datasets)
}

// ExperimentsHandler is a handler for '/experiments' endpoint.
// Lists all experiments in database (paginated only).
func ExperimentsHandler(c *gin.Context) {
	// Optional parameters (see CellsHandler)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	// Set max limit per_page to 1000
	if limit > 1000 {
		limit = 1000
	}
	experiments, err := PaginatedExperiments(page, limit)
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
	RenderJSONwithMeta(c, indent, page, limit, count, include, experiments)
}
