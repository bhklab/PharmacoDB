package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RootHandler handles GET requests for root endpoint.
func RootHandler(c *gin.Context) {
	u := "Welcome to PharmacoDB API.\n"
	v := "Current version: " + Version() + "\n"
	w := "Visit https://github.com/bhklab/PharmacoDB for more information."
	c.String(http.StatusOK, u+v+w)
}

// IndexCells returns a list of cell lines.
// Handles GET requests for /cell_lines.
func IndexCells(c *gin.Context) {
	var cells Cells
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	all, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	if all {
		err := cells.List()
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(cells)))
		RenderJSON(c, indent, cells)
	} else {
		err := cells.ListPaginated(page, limit)
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		total, err := Count("cells")
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		WriteHeader(c, "/cell_lines", page, limit, total)
		RenderJSONwithMeta(c, indent, page, limit, total, include, cells)
	}
}

// ShowCell returns a single cell line using id/name.
// Handles GET requests for /cell_lines/{id}.
func ShowCell(c *gin.Context) {
	var cell Cell
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	typ := c.DefaultQuery("type", "id")
	id := c.Param("id")
	err := cell.Find(id, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	err = cell.Annotate()
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, cell)
}

// CellCompounds returns all distinct compounds where a cell line of interest has been tested.
// Handles GET requests for /cell_lines/{id}/compounds.
func CellCompounds(c *gin.Context) {
	var cell Cell
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	err := cell.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	cellCompounds, total, err := cell.Compounds(page, limit)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	if total == 0 {
		NotFound(c, nil)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, cellCompounds)
}

// IndexTissues returns a list of tissues.
// Handles GET requests for /tissues.
func IndexTissues(c *gin.Context) {
	var tissues Tissues
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	if listAll {
		err := tissues.List()
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(tissues)))
		RenderJSON(c, indent, tissues)
	} else {
		err := tissues.ListPaginated(page, limit)
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		total, err := Count("tissues")
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		WriteHeader(c, "/tissues", page, limit, total)
		RenderJSONwithMeta(c, indent, page, limit, total, include, tissues)
	}
}

// ShowTissue returns a single tissue.
// Handles GET requests for /tissues/{id}.
func ShowTissue(c *gin.Context) {
	var tissue Tissue
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	typ := c.DefaultQuery("type", "id")
	id := c.Param("id")
	err := tissue.Find(id, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	err = tissue.Annotate()
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, tissue)
}

// TissueCells returns all cell lines of a tissue type.
// Handles GET requests for /tissues/{id}/cell_lines.
func TissueCells(c *gin.Context) {
	var tissue Tissue
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	err := tissue.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	tissueCells, total, err := tissue.Cells(page, limit)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	if total == 0 {
		NotFound(c, nil)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, tissueCells)
}

// TissueCompounds returns all compounds tested with tissue.
// Handles GET requests for /tissues/{id}/compounds.
func TissueCompounds(c *gin.Context) {
	var tissue Tissue
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	err := tissue.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	tissueCompounds, total, err := tissue.Compounds(page, limit)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	if total == 0 {
		NotFound(c, nil)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, tissueCompounds)
}

// IndexCompounds returns a list of compounds.
// Handles GET requests for /compounds.
func IndexCompounds(c *gin.Context) {
	var compounds Compounds
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	if listAll {
		err := compounds.List()
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(compounds)))
		RenderJSON(c, indent, compounds)
	} else {
		err := compounds.ListPaginated(page, limit)
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		total, err := Count("drugs")
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		WriteHeader(c, "/compounds", page, limit, total)
		RenderJSONwithMeta(c, indent, page, limit, total, include, compounds)
	}
}

// ShowCompound returns a single compound.
// Handles GET requests for /compounds/{id}.
func ShowCompound(c *gin.Context) {
	var compound Compound
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	typ := c.DefaultQuery("type", "id")
	id := c.Param("id")
	err := compound.Find(id, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	err = compound.Annotate()
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, compound)
}

// CompoundCells returns all cell lines tested with compound.
// Handles GET requests for /compounds/{id}/cell_lines.
func CompoundCells(c *gin.Context) {
	var compound Compound
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	err := compound.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	compoundCells, total, err := compound.Cells(page, limit)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	if total == 0 {
		NotFound(c, nil)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, compoundCells)
}

// CompoundTissues returns all tissues tested with compound.
// Handles GET requests for /compounds/{id}/tissues.
func CompoundTissues(c *gin.Context) {
	var compound Compound
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	err := compound.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	compoundTissues, total, err := compound.Tissues(page, limit)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	if total == 0 {
		NotFound(c, nil)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, compoundTissues)
}

// IndexDatasets returns a list of datasets.
// Handles GET requests for /datasets.
func IndexDatasets(c *gin.Context) {
	var datasets Datasets
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	listAll, _ := strconv.ParseBool(c.DefaultQuery("all", "false"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	if listAll {
		err := datasets.List()
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		c.Writer.Header().Set("Total-Records", strconv.Itoa(len(datasets)))
		RenderJSON(c, indent, datasets)
	} else {
		err := datasets.ListPaginated(page, limit)
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		total, err := Count("datasets")
		if err != nil {
			InternalServerError(c, nil)
			return
		}
		WriteHeader(c, "/datasets", page, limit, total)
		RenderJSONwithMeta(c, indent, page, limit, total, include, datasets)
	}
}

// ShowDataset returns a single dataset.
// Handles GET requests for /datasets/{id}.
func ShowDataset(c *gin.Context) {
	var dataset Dataset
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	typ := c.DefaultQuery("type", "id")
	id := c.Param("id")
	err := dataset.Find(id, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	RenderJSON(c, indent, dataset)
}

// DatasetCells returns a list of cell lines tested in dataset.
// Handles GET requests for /datasets/{id}/cell_lines.
func DatasetCells(c *gin.Context) {
	var dataset Dataset
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	err := dataset.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	datasetCells, total, err := dataset.Cells(page, limit)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	if total == 0 {
		NotFound(c, nil)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, datasetCells)
}

// DatasetTissues returns a list of tissues tested in dataset.
// Handles GET requests for /datasets/{id}/tissues.
func DatasetTissues(c *gin.Context) {
	var dataset Dataset
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	err := dataset.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	datasetTissues, total, err := dataset.Tissues(page, limit)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	if total == 0 {
		NotFound(c, nil)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, datasetTissues)
}

// DatasetCompounds returns a list of compounds tested in dataset.
// Handles GET requests for /datasets/{id}/compounds.
func DatasetCompounds(c *gin.Context) {
	var dataset Dataset
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	err := dataset.Find(c.Param("id"), c.DefaultQuery("type", "id"))
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	datasetCompounds, total, err := dataset.Compounds(page, limit)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	if total == 0 {
		NotFound(c, nil)
		return
	}
	RenderJSONwithMeta(c, indent, page, limit, total, include, datasetCompounds)
}

// IndexExperiments returns a list of experiments.
// Handles GET requests for /experiments.
func IndexExperiments(c *gin.Context) {
	var experiments Experiments
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	include := c.Query("include")
	// Set max limit per_page to 1000
	if limit > 1000 {
		limit = 1000
	}
	err := experiments.ListPaginated(page, limit)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	total, err := Count("experiments")
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	WriteHeader(c, "/experiments", page, limit, total)
	RenderJSONwithMeta(c, indent, page, limit, total, include, experiments)
}

// ShowExperiment returns a single experiment.
// Handles GET requests for /experiments/{id}.
func ShowExperiment(c *gin.Context) {
	var experiment Experiment
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	id := c.Param("id")
	err := experiment.Find(id)
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	err = experiment.DoseResponse()
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, experiment)
}

// IndexIntersections returns a list of all intersections.
// Handles GET requests for /intersections.
func IndexIntersections(c *gin.Context) {
	var intersections Intersections
	intersections.List()
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	RenderJSON(c, indent, intersections)
}

// CellCompoundIntersection returns a list of experiments where
// a cell line and a compound have been tested.
// Handles GET requests for /intersections/{id}/{cell_id}/{compound_id}.
func CellCompoundIntersection(c *gin.Context) {
	var experiments Experiments
	cellID := c.Param("cell_id")
	compoundID := c.Param("compound_id")
	typ := c.DefaultQuery("type", "id")
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	err := experiments.CellCompoundCombination(cellID, compoundID, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	RenderJSON(c, indent, experiments)
}

// CellDatasetIntersection returns a list of experiments where
// a cell line has been tested in a dataset.
// Handles GET requests for /intersections/{id}/{cell_id}/{dataset_id}.
func CellDatasetIntersection(c *gin.Context) {
	var experiments Experiments
	cellID := c.Param("cell_id")
	datasetID := c.Param("dataset_id")
	typ := c.DefaultQuery("type", "id")
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	err := experiments.CellDatasetCombination(cellID, datasetID, typ)
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(c, nil)
		} else {
			InternalServerError(c, nil)
		}
		return
	}
	RenderJSON(c, indent, experiments)
}

// StatTissuesCells returns a list of tissues,
// and the number of cell lines per tissue.
// Handles GET requests for /stats/tissues/cell_lines.
func StatTissuesCells(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	data, err := CountCellsPerTissue()
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}

// StatDatasetsCells returns a list of datasets,
// and the number of cell lines tested in each dataset.
// Handles GET requests for /stats/datasets/cell_lines.
func StatDatasetsCells(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	query := "SELECT dataset_id, dataset_name, cell_lines FROM source_statistics;"
	data, err := CountItemsPerDataset(query)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}

// StatDatasetsTissues returns a list of datasets,
// and the number of tissues tested in each dataset.
// Handles GET requests for /stats/datasets/tissues.
func StatDatasetsTissues(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	query := "SELECT dataset_id, dataset_name, tissues FROM source_statistics;"
	data, err := CountItemsPerDataset(query)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}

// StatDatasetsCellsCompounds returns a list of datasets,
// and the number of compounds tested with a cell line per dataset.
// Handles GET requests for /stats/datasets/compounds/cell_lines/:id.
func StatDatasetsCellsCompounds(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.drug_id) FROM experiments e WHERE e.cell_id = '%s' AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}

// StatDatasetsTissuesCells returns a list of datasets, and the
// number of cell lines in a tissue per dataset.
// Handles GET requests for /stats/datasets/cell_lines/tissues/{id}.
func StatDatasetsTissuesCells(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.cell_id) FROM experiments e WHERE e.tissue_id = '%s' AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}

// StatDatasetsTissuesCompounds returns a list of datasets, and the
// number of compounds tested with a tissue per dataset.
// Handles GET requests for /stats/datasets/compounds/tissues/{id}.
func StatDatasetsTissuesCompounds(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.drug_id) FROM experiments e WHERE e.tissue_id = '%s' AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}

// StatDatasetsCompounds returns a list of datasets, and the number of
// compounds tested in each dataset.
// Handles GET requests for /stats/datasets/compounds.
func StatDatasetsCompounds(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	query := "SELECT dataset_id, dataset_name, drugs FROM source_statistics;"
	data, err := CountItemsPerDataset(query)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}

// StatDatasetsCompoundsCells returns a list of datasets, and the number
// of cell lines tested with a drug per dataset.
// Handles GET requests for /stats/datasets/cell_lines/compounds/{id}.
func StatDatasetsCompoundsCells(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.cell_id) FROM experiments e WHERE e.drug_id = '%s' AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}

// StatDatasetsCompoundsTissues returns a list of datasets, and the
// number of tissues tested with a compound per dataset.
// Handles GET requests for /stats/datasets/tissues/compounds/{id}.
func StatDatasetsCompoundsTissues(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	id := c.Param("id")
	query := fmt.Sprintf("SELECT d.dataset_id, d.dataset_name, (SELECT COUNT(DISTINCT e.tissue_id) FROM experiments e WHERE e.drug_id = '%s' AND e.dataset_id = d.dataset_id) AS count FROM datasets d;", id)
	data, err := CountItemsPerDataset(query)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}

// StatDatasetsExperiments returns a list of datasets, and the number
// of experiments carried out in each dataset.
// Handles GET requests for /stats/datasets/experiments.
func StatDatasetsExperiments(c *gin.Context) {
	indent, _ := strconv.ParseBool(c.DefaultQuery("indent", "false"))
	query := "SELECT dataset_id, dataset_name, experiments FROM source_statistics;"
	data, err := CountItemsPerDataset(query)
	if err != nil {
		InternalServerError(c, nil)
		return
	}
	RenderJSON(c, indent, data)
}
