package api

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
		return -1, err
	}
	query := "SELECT COUNT(*) FROM " + table + ";"
	row := db.QueryRow(query)
	err = row.Scan(&count)
	if err != nil {
		LogPrivateError(err)
		return -1, err
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

// RenderJSONwithMeta outputs response as either indented or non-indented along with metadata about the reponse.
// Metadata includes current page, last page, per_page count and total count of result records.
func RenderJSONwithMeta(c *gin.Context, indent bool, page int, limit int, total int, include string, obj interface{}) {
	var data interface{}
	if include == "metadata" {
		lastPage := int(math.Ceil(float64(total) / float64(limit)))
		meta := gin.H{"page": page, "per_page": limit, "last_page": lastPage, "total": total}
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

// WriteHeader writes  metadata to response header.
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

// IndexCell is a handler for '/cell_lines' endpoint.
// Lists all cell lines in database (paginated or non-paginated).
func IndexCell(c *gin.Context) {
	var cells Cells
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
		err := cells.List()
		if err != nil {
			LogInternalServerError(c)
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(cells)))
		RenderJSON(c, indent, cells)
		return
	}
	err := cells.ListPaginated(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	total, err := Count("cells")
	if err != nil {
		LogInternalServerError(c)
		return
	}
	WriteHeader(c, "/cell_lines", page, limit, total)
	RenderJSONwithMeta(c, indent, page, limit, total, include, cells)
}

// ShowCell is a handler for '/cell_lines/:id' endpoint.
// Returns a single cell line.
func ShowCell(c *gin.Context) {
	var cell Cell
	// Optional parameters
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	typ := c.DefaultQuery("type", "id")
	id := c.Param("id")
	err := cell.Find(id, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	err = cell.Annotate()
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, cell)
}

// IndexTissue is a handler for '/tissues' endpoint.
// Lists all tissues in database (paginated or non-paginated).
func IndexTissue(c *gin.Context) {
	var tissues Tissues
	// Optional parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	if listAll {
		err := tissues.List()
		if err != nil {
			LogInternalServerError(c)
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(tissues)))
		RenderJSON(c, indent, tissues)
		return
	}
	err := tissues.ListPaginated(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	total, err := Count("tissues")
	if err != nil {
		LogInternalServerError(c)
		return
	}
	WriteHeader(c, "/tissues", page, limit, total)
	RenderJSONwithMeta(c, indent, page, limit, total, include, tissues)
}

// ShowTissue is a handler for '/tissues/:id' endpoint.
// Returns a single tissue.
func ShowTissue(c *gin.Context) {
	var tissue Tissue
	// Optional parameters
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	typ := c.DefaultQuery("type", "id")
	id := c.Param("id")
	err := tissue.Find(id, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	err = tissue.Annotate()
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, tissue)
}

// IndexDrug is a handler for '/drugs' endpoint.
// Lists all drugs in database (paginated or non-paginated).
func IndexDrug(c *gin.Context) {
	var drugs Drugs
	// Optional parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	if listAll {
		err := drugs.List()
		if err != nil {
			LogInternalServerError(c)
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(drugs)))
		RenderJSON(c, indent, drugs)
		return
	}
	err := drugs.ListPaginated(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	total, err := Count("drugs")
	if err != nil {
		LogInternalServerError(c)
		return
	}
	WriteHeader(c, "/drugs", page, limit, total)
	RenderJSONwithMeta(c, indent, page, limit, total, include, drugs)
}

// ShowDrug is a handler for '/drugs/:id' endpoint.
// Returns a single drug.
func ShowDrug(c *gin.Context) {
	var drug Drug
	// Optional parameters
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	typ := c.DefaultQuery("type", "id")
	id := c.Param("id")
	err := drug.Find(id, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	err = drug.Annotate()
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, drug)
}
