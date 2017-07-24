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

// CellDrugs is a handler for '/cell_lines/:id/drugs' endpoint.
// Lists all distinct drugs where a cell line of interest has been tested,
// along with datasets and experiment count.
func CellDrugs(c *gin.Context) {
	var cell Cell
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	err := cell.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	celldrugs, total, err := cell.Drugs(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	if total == 0 {
		LogNotFoundError(c)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, celldrugs)
}

// IndexTissue is a handler for '/tissues' endpoint.
// Lists all tissues in database (paginated or non-paginated).
func IndexTissue(c *gin.Context) {
	var tissues Tissues
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

// TissueCells is a handler for '/tissues/:id/cell_lines' endpoint.
// Lists all cell lines of a certain tissue type.
func TissueCells(c *gin.Context) {
	var tissue Tissue
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	err := tissue.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	tissueCells, total, err := tissue.Cells(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	if total == 0 {
		LogNotFoundError(c)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, tissueCells)
}

// TissueDrugs is a handler for '/tissues/:id/drugs' endpoint.
// Lists all distinct drugs where a tissue of interest has been tested,
// along with datasets and experiment count.
func TissueDrugs(c *gin.Context) {
	var tissue Tissue
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	err := tissue.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	tissueDrugs, total, err := tissue.Drugs(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	if total == 0 {
		LogNotFoundError(c)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, tissueDrugs)
}

// IndexDrug is a handler for '/drugs' endpoint.
// Lists all drugs in database (paginated or non-paginated).
func IndexDrug(c *gin.Context) {
	var drugs Drugs
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

// DrugCells is a handler for '/drugs/:id/cell_lines' endpoint.
// Lists all distinct cell lines where a drug of interest has been tested,
// along with datasets and experiment count.
func DrugCells(c *gin.Context) {
	var drug Drug
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	err := drug.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	drugCells, total, err := drug.Cells(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	if total == 0 {
		LogNotFoundError(c)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, drugCells)
}

// IndexDataset is a handler for '/datasets' endpoint.
// Lists all datasets in database (paginated or non-paginated).
func IndexDataset(c *gin.Context) {
	var datasets Datasets
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	if listAll {
		err := datasets.List()
		if err != nil {
			LogInternalServerError(c)
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(datasets)))
		RenderJSON(c, indent, datasets)
		return
	}
	err := datasets.ListPaginated(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	total, err := Count("datasets")
	if err != nil {
		LogInternalServerError(c)
		return
	}
	WriteHeader(c, "/datasets", page, limit, total)
	RenderJSONwithMeta(c, indent, page, limit, total, include, datasets)
}

// ShowDataset is a handler for '/datasets/:id' endpoint.
// Returns a single dataset.
func ShowDataset(c *gin.Context) {
	var dataset Dataset
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	typ := c.DefaultQuery("type", "id")
	id := c.Param("id")
	err := dataset.Find(id, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	RenderJSON(c, indent, dataset)
}

// IndexExperiment is a handler for '/experiments' endpoint.
// Lists all experiments in database, with pagination only.
func IndexExperiment(c *gin.Context) {
	var experiments Experiments
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	// Set max limit per_page to 1000
	if limit > 1000 {
		limit = 1000
	}
	err := experiments.ListPaginated(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	total, err := Count("experiments")
	if err != nil {
		LogInternalServerError(c)
		return
	}
	WriteHeader(c, "/experiments", page, limit, total)
	RenderJSONwithMeta(c, indent, page, limit, total, include, experiments)
}

// ShowExperiment is a handler for '/experiments/i/:id' endpoint.
// Returns a single experiment with associated dose/response data.
func ShowExperiment(c *gin.Context) {
	var experiment Experiment
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	id := c.Param("id")
	err := experiment.Find(id)
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	err = experiment.DoseResponse()
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, experiment)
}

// CellDrugExperiments is a handler for '/experiments/x/:cell_id/:drug_id' endpoint.
// Lists all experiments (including dose/response data) for a cell line and drug combination.
func CellDrugExperiments(c *gin.Context) {
	var experiments Experiments
	cellID := c.Param("cell_id")
	drugID := c.Param("drug_id")
	typ := c.DefaultQuery("type", "id")
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	err := experiments.CellDrugCombination(cellID, drugID, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	RenderJSON(c, indent, experiments)
}

// CellDatasetExperiments is a handler for '/experiments/y/:cell_id/:dataset_id' endpoint.
// Lists all experiments (including dose/response data) for a cell line and dataset combination.
func CellDatasetExperiments(c *gin.Context) {
	var experiments Experiments
	cellID := c.Param("cell_id")
	datasetID := c.Param("dataset_id")
	typ := c.DefaultQuery("type", "id")
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	err := experiments.CellDatasetCombination(cellID, datasetID, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	RenderJSON(c, indent, experiments)
}
