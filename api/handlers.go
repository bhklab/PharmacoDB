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

// RenderJSONwithMeta outputs response as either indented or non-indented along with metadata about the response.
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

// DrugTissues is a handler for '/drugs/:id/tissues' endpoint.
// Lists all distinct tissues where a drug of interest has been tested,
// along with datasets and experiment count.
func DrugTissues(c *gin.Context) {
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
	drugTissues, total, err := drug.Tissues(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	if total == 0 {
		LogNotFoundError(c)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, drugTissues)
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

// DatasetCells is a handler for '/datasets/:id/cell_lines' endpoint.
// Lists all distinct cell lines which have been tested in a dataset of interest.
func DatasetCells(c *gin.Context) {
	var dataset Dataset
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	err := dataset.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	datasetCells, total, err := dataset.Cells(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	if total == 0 {
		LogNotFoundError(c)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, datasetCells)
}

// DatasetTissues is a handler for '/datasets/:id/tissues' endpoint.
// Lists all distinct tissues that have been tested with a dataset of interest.
func DatasetTissues(c *gin.Context) {
	var dataset Dataset
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	err := dataset.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	datasetTissues, total, err := dataset.Tissues(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	if total == 0 {
		LogNotFoundError(c)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, datasetTissues)
}

// DatasetDrugs is a handler for '/datasets/:id/drugs' endpoint.
// Lists all distinct drugs which have been tested in a dataset of interest.
func DatasetDrugs(c *gin.Context) {
	var dataset Dataset
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	include := c.Query("include")
	err := dataset.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			LogNotFoundError(c)
		} else {
			LogInternalServerError(c)
		}
		return
	}
	datasetDrugs, total, err := dataset.Drugs(page, limit)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	if total == 0 {
		LogNotFoundError(c)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, datasetDrugs)
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

// ShowExperiment is a handler for '/experiments/:id' endpoint.
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

// IndexIntersection is a handler for '/intersections' endpoint.
// Lists all available intersections in API.
func IndexIntersection(c *gin.Context) {
	var intersections Intersections
	intersections.List()
	RenderJSON(c, true, intersections)
}

// CellDrugIntersection is a handler for '/intersections/1/:cell_id/:drug_id' endpoint.
// Lists all experiments (including dose/response data) for a cell line and drug combination.
func CellDrugIntersection(c *gin.Context) {
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

// CellDatasetIntersection is a handler for '/intersections/2/:cell_id/:dataset_id' endpoint.
// Lists all experiments (including dose/response data) for a cell line and dataset combination.
func CellDatasetIntersection(c *gin.Context) {
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

// TissueCellStats is a handler for '/stats/tissues/cell_lines' endpoint.
// Lists all tissues, along with the number of cell lines in each tissue.
func TissueCellStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	data, err := CountCellsPerTissue()
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}

// DatasetCellStats is a handler for '/stats/datasets/cell_lines' endpoint.
// Lists all datasets, along with the number of cell lines tested in each dataset.
func DatasetCellStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	query := "SELECT dataset_id, dataset_name, cell_lines FROM source_statistics;"
	data, err := CountItemsPerDataset(query)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}

// DatasetCellDrugsStats is a handler for '/stats/datasets/cell_lines/:id/drugs' endpoint.
// Lists all datasets, along with the number of drugs tested with a cell line per dataset.
func DatasetCellDrugsStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.drug_id) FROM experiments e WHERE e.cell_id = %s AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}

// DatasetTissueStats is a handler for '/stats/datasets/tissues' endpoint.
// Lists all datasets, along with the number of tissues tested in each dataset.
func DatasetTissueStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	query := "SELECT dataset_id, dataset_name, tissues FROM source_statistics;"
	data, err := CountItemsPerDataset(query)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}

// DatasetTissueCellsStats is a handler for '/stats/datasets/tissues/:id/cell_lines' endpoint.
// Lists all datasets, along with the number of cell_lines in a tissue per dataset.
func DatasetTissueCellsStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.cell_id) FROM experiments e WHERE e.tissue_id = %s AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}

// DatasetTissueDrugsStats is a handler for '/stats/datasets/tissues/:id/drugs' endpoint.
// Lists all datasets, along with the number of drugs tested with a tissue per dataset.
func DatasetTissueDrugsStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.drug_id) FROM experiments e WHERE e.tissue_id = %s AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}

// DatasetDrugStats is a handler for '/stats/datasets/drugs' endpoint.
// Lists all datasets, along with the number of drugs tested in each dataset.
func DatasetDrugStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	query := "SELECT dataset_id, dataset_name, drugs FROM source_statistics;"
	data, err := CountItemsPerDataset(query)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}

// DatasetDrugCellsStats is a handler for '/stats/datasets/drugs/:id/cell_lines' endpoint.
// Lists all datasets, along with the number of cell_lines tested with drug per dataset.
func DatasetDrugCellsStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.cell_id) FROM experiments e WHERE e.drug_id = %s AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}

// DatasetDrugTissuesStats is a handler for '/stats/datasets/drugs/:id/tissues' endpoint.
// Lists all datasets, along with the number of tissues tested with drug per dataset.
func DatasetDrugTissuesStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.tissue_id) FROM experiments e WHERE e.drug_id = %s AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}

// DatasetExperimentStats is a handler for '/stats/datasets/experiments' endpoint.
// Lists all datasets, along with the number of experiments tested in each dataset.
func DatasetExperimentStats(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "true"))
	query := "SELECT dataset_id, dataset_name, experiments FROM source_statistics;"
	data, err := CountItemsPerDataset(query)
	if err != nil {
		LogInternalServerError(c)
		return
	}
	RenderJSON(c, indent, data)
}
