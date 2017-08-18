package api

import (
	"database/sql"
	"fmt"
	"strings"
)

// Comp models experiment data for
// a compound-resource pair test.
type Comp struct {
	Compound Compound `json:"compound"`
	Datasets []string `json:"datasets"`
	Count    int      `json:"experiment_count"`
}

// Comps is a collection of Comp.
type Comps []Comp

// DC models experiment data for
// a compound and a cell line.
type DC struct {
	Cell     Cell     `json:"cell_line"`
	Datasets []string `json:"datasets"`
	Count    int      `json:"experiment_count"`
}

// DCS is a collection of DCT.
type DCS []DC

// DT models experiment data for
// a compound and a tissue.
type DT struct {
	Tissue   Tissue   `json:"tissue"`
	Datasets []string `json:"datasets"`
	Count    int      `json:"experiment_count"`
}

// DTS is a collection of DT.
type DTS []DT

// Intersection is a combination query model.
type Intersection struct {
	ID      int      `json:"id"`
	ResComb []string `json:"resources_combined"`
	Path    string   `json:"template_path"`
	Ex      string   `json:"example_path"`
}

// Intersections is a collection of Intersection.
type Intersections []Intersection

// TissueCount models the number of cell lines per tissue.
type TissueCount struct {
	Tissue Tissue `json:"tissue"`
	Count  int    `json:"cell_lines_count"`
}

// TissueCounts is a collection of TissueCount.
type TissueCounts []TissueCount

// DatasetCount models the number of compounds tested per dataset.
type DatasetCount struct {
	Dataset Dataset `json:"dataset"`
	Count   int     `json:"compounds_count"`
}

// DatasetCounts is a collection of DatasetCount.
type DatasetCounts []DatasetCount

// List returns a list of all cell lines, without pagination.
func (cells *Cells) List() error {
	var cell Cell
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT cell_id, cell_name FROM cells;")
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		*cells = append(*cells, cell)
	}
	return nil
}

// ListPaginated returns a list of cell lines with pagination.
func (cells *Cells) ListPaginated(page int, limit int) error {
	var cell Cell
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT cell_id, cell_name FROM cells LIMIT %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		*cells = append(*cells, cell)
	}
	return nil
}

// Find updates receiver with a record for a single cell line.
func (cell *Cell) Find(id string, typ string) error {
	var query string
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	if isSameStr(typ, "name") {
		query = "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE c.cell_name LIKE ?;"
	} else if isSameStr(typ, "accession") {
		query = "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE c.accession_id LIKE ?;"
	} else {
		query = "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE c.cell_id LIKE ?;"
	}
	tissue := &Tissue{}
	row := db.QueryRow(query, id)
	err = row.Scan(&cell.ID, &cell.ACC, &cell.Name, &tissue.ID, &tissue.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogSentry(err)
		}
		return err
	}
	cell.Tissue = tissue
	return nil
}

// Annotate adds annotations to cell.
func (cell *Cell) Annotate() error {
	var (
		annotation     Annotation
		annotations    Annotations
		annotationName string
		datasetName    string
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	query := "SELECT s.cell_name, d.dataset_name FROM source_cell_names s JOIN datasets d ON d.dataset_id = s.source_id WHERE s.cell_id = ?;"
	rows, err := db.Query(query, cell.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	exists := make(map[string]bool)
	for rows.Next() {
		err = rows.Scan(&annotationName, &datasetName)
		if err != nil {
			LogSentry(err)
			return err
		}
		if exists[annotationName] {
			for i, a := range annotations {
				if a.Name == annotationName && !stringInSlice(datasetName, a.Datasets) {
					annotations[i].Datasets = append(annotations[i].Datasets, datasetName)
				}
			}
		} else {
			var datasetsNew []string
			annotation.Name = annotationName
			annotation.Datasets = append(datasetsNew, datasetName)
			annotations = append(annotations, annotation)
			exists[annotationName] = true
		}
	}
	cell.Annotations = annotations
	return nil
}

// Compounds updates cell with a list of distinct compounds tested with cell line.
func (cell *Cell) Compounds(page int, limit int) (Comps, int, error) {
	var (
		comp     Comp
		comps    Comps
		datasets string
		count    int
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return comps, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS d.drug_id, d.drug_name, GROUP_CONCAT(DISTINCT da.dataset_name) AS datasets, COUNT(*) AS experiment_count FROM experiments e JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.cell_id = ? GROUP BY e.drug_id ORDER BY experiment_count DESC LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, cell.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return comps, count, err
	}
	for rows.Next() {
		err = rows.Scan(&comp.Compound.ID, &comp.Compound.Name, &datasets, &comp.Count)
		if err != nil {
			LogSentry(err)
			return comps, count, err
		}
		comp.Datasets = strings.Split(datasets, ",")
		comps = append(comps, comp)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogSentry(err)
		return comps, count, err
	}
	return comps, count, nil
}

// List updates receiver with a list of all tissues without pagination.
func (tissues *Tissues) List() error {
	var tissue Tissue
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT tissue_id, tissue_name FROM tissues;")
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		*tissues = append(*tissues, tissue)
	}
	return nil
}

// ListPaginated updates receiver with a list of tissues using pagination.
func (tissues *Tissues) ListPaginated(page int, limit int) error {
	var tissue Tissue
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT tissue_id, tissue_name FROM tissues LIMIT %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		*tissues = append(*tissues, tissue)
	}
	return nil
}

// Find updates receiver with a record for a single tissue.
func (tissue *Tissue) Find(id string, typ string) error {
	var query string
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	if isSameStr(typ, "name") {
		query = "SELECT tissue_id, tissue_name FROM tissues WHERE tissue_name LIKE ?;"
	} else {
		query = "SELECT tissue_id, tissue_name FROM tissues WHERE tissue_id LIKE ?;"
	}
	row := db.QueryRow(query, id)
	err = row.Scan(&tissue.ID, &tissue.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogSentry(err)
		}
		return err
	}
	return nil
}

// Annotate adds annotations to tissue.
func (tissue *Tissue) Annotate() error {
	var (
		annotation     Annotation
		annotations    Annotations
		annotationName string
		datasetName    string
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	query := "SELECT s.tissue_name, d.dataset_name FROM source_tissue_names s JOIN datasets d ON d.dataset_id = s.source_id WHERE s.tissue_id = ?;"
	rows, err := db.Query(query, tissue.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	exists := make(map[string]bool)
	for rows.Next() {
		err = rows.Scan(&annotationName, &datasetName)
		if err != nil {
			LogSentry(err)
			return err
		}
		if exists[annotationName] {
			for i, a := range annotations {
				if a.Name == annotationName && !stringInSlice(datasetName, a.Datasets) {
					annotations[i].Datasets = append(annotations[i].Datasets, datasetName)
				}
			}
		} else {
			var datasetsNew []string
			annotation.Name = annotationName
			annotation.Datasets = append(datasetsNew, datasetName)
			annotations = append(annotations, annotation)
			exists[annotationName] = true
		}
	}
	tissue.Annotations = annotations
	return nil
}

// Cells returns a paginated list of all cell lines of tissue type.
func (tissue *Tissue) Cells(page int, limit int) (Cells, int, error) {
	var (
		cell  Cell
		cells Cells
		count int
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return cells, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS cell_id, cell_name FROM cells WHERE tissue_id = ? LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, tissue.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return cells, count, err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			LogSentry(err)
			return cells, count, err
		}
		cells = append(cells, cell)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogSentry(err)
		return cells, count, err
	}
	return cells, count, nil
}

// Compounds returns a paginated list of all distinct compounds that have been tested with tissue.
func (tissue *Tissue) Compounds(page int, limit int) (Comps, int, error) {
	var (
		comp     Comp
		comps    Comps
		datasets string
		count    int
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return comps, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS d.drug_id, d.drug_name, GROUP_CONCAT(DISTINCT da.dataset_name) AS datasets, COUNT(*) AS experiment_count FROM experiments e JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.tissue_id = ? GROUP BY e.drug_id ORDER BY experiment_count DESC LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, tissue.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return comps, count, err
	}
	for rows.Next() {
		err = rows.Scan(&comp.Compound.ID, &comp.Compound.Name, &datasets, &comp.Count)
		if err != nil {
			LogSentry(err)
			return comps, count, err
		}
		comp.Datasets = strings.Split(datasets, ",")
		comps = append(comps, comp)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogSentry(err)
		return comps, count, err
	}
	return comps, count, nil
}

// List updates receiver with a list of all compounds without pagination.
func (compounds *Compounds) List() error {
	var compound Compound
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT drug_id, drug_name FROM drugs;")
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&compound.ID, &compound.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		*compounds = append(*compounds, compound)
	}
	return nil
}

// ListPaginated updates receiver with a list of compounds using pagination.
func (compounds *Compounds) ListPaginated(page int, limit int) error {
	var compound Compound
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT drug_id, drug_name FROM drugs LIMIT %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&compound.ID, &compound.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		*compounds = append(*compounds, compound)
	}
	return nil
}

// Find updates receiver with a record for a single compound.
func (compound *Compound) Find(id string, typ string) error {
	var query string
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	if isSameStr(typ, "name") {
		query = "SELECT drug_id, drug_name FROM drugs WHERE drug_name LIKE ?;"
	} else {
		query = "SELECT drug_id, drug_name FROM drugs WHERE drug_id LIKE ?;"
	}
	row := db.QueryRow(query, id)
	err = row.Scan(&compound.ID, &compound.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogSentry(err)
		}
		return err
	}
	return nil
}

// Annotate adds annotations to compound.
func (compound *Compound) Annotate() error {
	var (
		annotation     Annotation
		annotations    Annotations
		annotationName string
		datasetName    string
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	query := "SELECT s.drug_name, d.dataset_name FROM source_drug_names s JOIN datasets d ON d.dataset_id = s.source_id WHERE s.drug_id = ?;"
	rows, err := db.Query(query, compound.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	exists := make(map[string]bool)
	for rows.Next() {
		err = rows.Scan(&annotationName, &datasetName)
		if err != nil {
			LogSentry(err)
			return err
		}
		if exists[annotationName] {
			for i, a := range annotations {
				if a.Name == annotationName && !stringInSlice(datasetName, a.Datasets) {
					annotations[i].Datasets = append(annotations[i].Datasets, datasetName)
				}
			}
		} else {
			var datasetsNew []string
			annotation.Name = annotationName
			annotation.Datasets = append(datasetsNew, datasetName)
			annotations = append(annotations, annotation)
			exists[annotationName] = true
		}
	}
	compound.Annotations = annotations
	return nil
}

// Cells returns a paginated list of all distinct cell lines which have been tested with drug, along with
// experiments count and an array of datasets that tested each cell/drug combination.
func (compound *Compound) Cells(page int, limit int) (DCS, int, error) {
	var (
		compoundCell  DC
		compoundCells DCS
		datasets      string
		count         int
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return compoundCells, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS c.cell_id, c.cell_name, GROUP_CONCAT(DISTINCT da.dataset_name) AS datasets, COUNT(*) AS experiment_count FROM experiments e JOIN cells c ON c.cell_id = e.cell_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.drug_id = ? GROUP BY e.cell_id ORDER BY experiment_count DESC LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, compound.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return compoundCells, count, err
	}
	for rows.Next() {
		err = rows.Scan(&compoundCell.Cell.ID, &compoundCell.Cell.Name, &datasets, &compoundCell.Count)
		if err != nil {
			LogSentry(err)
			return compoundCells, count, err
		}
		compoundCell.Datasets = strings.Split(datasets, ",")
		compoundCells = append(compoundCells, compoundCell)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogSentry(err)
		return compoundCells, count, err
	}
	return compoundCells, count, nil
}

// Tissues returns a paginated list of all distinct tissues which have been tested with drug, along with
// experiments count and an array of datasets that tested each cell/drug combination.
func (compound *Compound) Tissues(page int, limit int) (DTS, int, error) {
	var (
		compoundTissue  DT
		compoundTissues DTS
		datasets        string
		count           int
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return compoundTissues, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS t.tissue_id, t.tissue_name, GROUP_CONCAT(DISTINCT da.dataset_name) AS datasets, COUNT(*) AS experiment_count FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.drug_id = ? GROUP BY e.tissue_id ORDER BY experiment_count DESC LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, compound.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return compoundTissues, count, err
	}
	for rows.Next() {
		err = rows.Scan(&compoundTissue.Tissue.ID, &compoundTissue.Tissue.Name, &datasets, &compoundTissue.Count)
		if err != nil {
			LogSentry(err)
			return compoundTissues, count, err
		}
		compoundTissue.Datasets = strings.Split(datasets, ",")
		compoundTissues = append(compoundTissues, compoundTissue)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogSentry(err)
		return compoundTissues, count, err
	}
	return compoundTissues, count, nil
}

// List updates receiver with a list of all datasets without pagination.
func (datasets *Datasets) List() error {
	var dataset Dataset
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT dataset_id, dataset_name FROM datasets;")
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&dataset.ID, &dataset.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		*datasets = append(*datasets, dataset)
	}
	return nil
}

// ListPaginated updates receiver with a list of datasets using pagination.
func (datasets *Datasets) ListPaginated(page int, limit int) error {
	var dataset Dataset
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT dataset_id, dataset_name FROM datasets LIMIT %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&dataset.ID, &dataset.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		*datasets = append(*datasets, dataset)
	}
	return nil
}

// Find updates receiver with a record for a single dataset.
func (dataset *Dataset) Find(id string, typ string) error {
	var query string
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	if isSameStr(typ, "name") {
		query = "SELECT dataset_id, dataset_name FROM datasets WHERE dataset_name LIKE ?;"
	} else {
		query = "SELECT dataset_id, dataset_name FROM datasets WHERE dataset_id LIKE ?;"
	}
	row := db.QueryRow(query, id)
	err = row.Scan(&dataset.ID, &dataset.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogSentry(err)
		}
		return err
	}
	return nil
}

// Cells returns a paginated list of all cell lines tested in a dataset.
func (dataset *Dataset) Cells(page int, limit int) (Cells, int, error) {
	var (
		cell  Cell
		cells Cells
		count int
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return cells, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS c.cell_id, c.cell_name FROM experiments e JOIN cells c ON c.cell_id = e.cell_id WHERE e.dataset_id = ? GROUP BY(e.cell_id) LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, dataset.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return cells, count, err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			LogSentry(err)
			return cells, count, err
		}
		cells = append(cells, cell)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogSentry(err)
		return cells, count, err
	}
	return cells, count, nil
}

// Tissues returns a paginated list of all tissues tested in dataset.
func (dataset *Dataset) Tissues(page int, limit int) (Tissues, int, error) {
	var (
		tissue  Tissue
		tissues Tissues
		count   int
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return tissues, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS t.tissue_id, t.tissue_name FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id WHERE e.dataset_id = ? GROUP BY(e.tissue_id) LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, dataset.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return tissues, count, err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogSentry(err)
			return tissues, count, err
		}
		tissues = append(tissues, tissue)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogSentry(err)
		return tissues, count, err
	}
	return tissues, count, nil
}

// Compounds returns a paginated list of all drugs tested in dataset.
func (dataset *Dataset) Compounds(page int, limit int) (Compounds, int, error) {
	var (
		compound  Compound
		compounds Compounds
		count     int
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return compounds, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS d.drug_id, d.drug_name FROM experiments e JOIN drugs d ON d.drug_id = e.drug_id WHERE e.dataset_id = ? GROUP BY(e.drug_id) LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, &dataset.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return compounds, count, err
	}
	for rows.Next() {
		err = rows.Scan(&compound.ID, &compound.Name)
		if err != nil {
			LogSentry(err)
			return compounds, count, err
		}
		compounds = append(compounds, compound)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogSentry(err)
		return compounds, count, err
	}
	return compounds, count, nil
}

// ListPaginated updates receiver with a list of experiments using pagination.
func (experiments *Experiments) ListPaginated(page int, limit int) error {
	var experiment Experiment
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id limit %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Compound.ID, &experiment.Compound.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		*experiments = append(*experiments, experiment)
	}
	return nil
}

// Find updates receiver with a record for a single experiment.
func (experiment *Experiment) Find(id string) error {
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	query := "SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.experiment_id = ?;"
	row := db.QueryRow(query, id)
	err = row.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Compound.ID, &experiment.Compound.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogSentry(err)
		}
		return err
	}
	return nil
}

// DoseResponse updates receiver experiment with its corresponding dose/response data.
func (experiment *Experiment) DoseResponse() error {
	var doseResponse DoseResponse
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT dose, response FROM dose_responses WHERE experiment_id = ?;", experiment.ID)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&doseResponse.Dose, &doseResponse.Response)
		if err != nil {
			LogSentry(err)
			return err
		}
		experiment.DoseResponses = append(experiment.DoseResponses, doseResponse)
	}
	return nil
}

// List lists all possible intersections.
func (intersections *Intersections) List() {
	var intersection Intersection
	// First intersection
	intersection.ID = 1
	intersection.ResComb = []string{"cell_line", "drug"}
	intersection.Path = "/intersections/{id}/{cell_id}/{drug_id}"
	intersection.Ex = "/intersections/1/mcf7/paclitaxel?type=name"
	*intersections = append(*intersections, intersection)
	// Second intersection
	intersection.ID = 2
	intersection.ResComb = []string{"cell_line", "dataset"}
	intersection.Path = "/intersections/{id}/{cell_id}/{dataset_id}"
	intersection.Ex = "/intersections/2/mcf7/ccle?type=name"
	*intersections = append(*intersections, intersection)
}

// CellCompoundCombination updates receiver with a list of all experiments where a cell line and a drug have been tested.
func (experiments *Experiments) CellCompoundCombination(cellID string, compoundID string, typ string) error {
	var (
		cell     Cell
		compound Compound
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	err = cell.Find(cellID, typ)
	if err != nil {
		return err
	}
	err = compound.Find(compoundID, typ)
	if err != nil {
		return err
	}
	query := "SELECT e.experiment_id, t.tissue_id, t.tissue_name, da.dataset_id, da.dataset_name FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.cell_id = ? AND e.drug_id = ?;"
	rows, _ := db.Query(query, cell.ID, compound.ID)
	defer rows.Close()
	for rows.Next() {
		var experiment Experiment
		err = rows.Scan(&experiment.ID, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		experiment.Cell.ID = cell.ID
		experiment.Cell.Name = cell.Name
		experiment.Compound = compound
		err = experiment.DoseResponse()
		if err != nil {
			LogSentry(err)
			return err
		}
		*experiments = append(*experiments, experiment)
	}
	return nil
}

// CellDatasetCombination updates receiver with a list of all experiments where a cell line and dataset have been tested.
func (experiments *Experiments) CellDatasetCombination(cellID string, datasetID string, typ string) error {
	var (
		cell    Cell
		dataset Dataset
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return err
	}
	err = cell.Find(cellID, typ)
	if err != nil {
		return err
	}
	err = dataset.Find(datasetID, typ)
	if err != nil {
		return err
	}
	query := "SELECT e.experiment_id, t.tissue_id, t.tissue_name, d.drug_id, d.drug_name FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id JOIN drugs d ON d.drug_id = e.drug_id WHERE e.cell_id = ? AND e.dataset_id = ?;"
	rows, _ := db.Query(query, cell.ID, dataset.ID)
	defer rows.Close()
	for rows.Next() {
		var experiment Experiment
		err = rows.Scan(&experiment.ID, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Compound.ID, &experiment.Compound.Name)
		if err != nil {
			LogSentry(err)
			return err
		}
		experiment.Cell.ID = cell.ID
		experiment.Cell.Name = cell.Name
		experiment.Dataset = dataset
		err = experiment.DoseResponse()
		if err != nil {
			LogSentry(err)
			return err
		}
		*experiments = append(*experiments, experiment)
	}
	return nil
}

// CountCellsPerTissue returns a list of all tissues, along with the number of
// cell lines of each tissue type.
func CountCellsPerTissue() (TissueCounts, error) {
	var (
		tissueCellCount  TissueCount
		tissueCellCounts TissueCounts
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return tissueCellCounts, err
	}
	query := "SELECT t.tissue_id, t.tissue_name, COUNT(*) AS cell_lines_count FROM tissues t JOIN cells c ON c.tissue_id = t.tissue_id GROUP BY(c.tissue_id);"
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return tissueCellCounts, err
	}
	for rows.Next() {
		err = rows.Scan(&tissueCellCount.Tissue.ID, &tissueCellCount.Tissue.Name, &tissueCellCount.Count)
		if err != nil {
			LogSentry(err)
			return tissueCellCounts, err
		}
		tissueCellCounts = append(tissueCellCounts, tissueCellCount)
	}
	return tissueCellCounts, nil
}

// CountItemsPerDataset returns a list of all datasets, along with the number of
// required item tested in each dataset.
func CountItemsPerDataset(query string) (DatasetCounts, error) {
	var (
		count  DatasetCount
		counts DatasetCounts
	)
	db, err := Database()
	defer db.Close()
	if err != nil {
		return counts, err
	}
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogSentry(err)
		return counts, err
	}
	for rows.Next() {
		err = rows.Scan(&count.Dataset.ID, &count.Dataset.Name, &count.Count)
		if err != nil {
			LogSentry(err)
			return counts, err
		}
		counts = append(counts, count)
	}
	return counts, nil
}
