package api

import (
	"database/sql"
	"fmt"
	"strings"
)

// CD models experiment data for
// a cell line and a drug.
type CD struct {
	Drug     Drug     `json:"drug"`
	Datasets []string `json:"datasets"`
	Count    int      `json:"experiment_count"`
}

// CDS is a collection of CD.
type CDS []CD

// List updates receiver with a list of all cell lines without pagination.
func (cells *Cells) List() error {
	var cell Cell
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT cell_id, accession_id, cell_name FROM cells;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*cells = append(*cells, cell)
	}
	return nil
}

// ListPaginated updates receiver with a list of cell lines using pagination.
func (cells *Cells) ListPaginated(page int, limit int) error {
	var cell Cell
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT cell_id, accession_id, cell_name FROM cells LIMIT %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*cells = append(*cells, cell)
	}
	return nil
}

// Find updates receiver with a record for a single cell line.
func (cell *Cell) Find(id string, typ string) error {
	var query string
	db, err := InitDB()
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
			LogPrivateError(err)
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
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	query := "SELECT s.cell_name, d.dataset_name FROM source_cell_names s JOIN datasets d ON d.dataset_id = s.source_id WHERE s.cell_id = ?;"
	rows, err := db.Query(query, cell.ID)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	exists := make(map[string]bool)
	for rows.Next() {
		err = rows.Scan(&annotationName, &datasetName)
		if err != nil {
			LogPrivateError(err)
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

// Drugs returns a paginated list of all distinct drugs which have been tested with cell, along with
// experiments count and an array of datasets that tested each cell/drug combination.
func (cell *Cell) Drugs(page int, limit int) (CDS, int, error) {
	var (
		cellDrug  CD
		cellDrugs CDS
		datasets  string
		count     int
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return cellDrugs, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS d.drug_id, d.drug_name, GROUP_CONCAT(DISTINCT da.dataset_name) AS datasets, COUNT(*) AS experiment_count FROM experiments e JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.cell_id = ? GROUP BY e.drug_id ORDER BY experiment_count DESC LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, cell.ID)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return cellDrugs, count, err
	}
	for rows.Next() {
		err = rows.Scan(&cellDrug.Drug.ID, &cellDrug.Drug.Name, &datasets, &cellDrug.Count)
		if err != nil {
			LogPrivateError(err)
			return cellDrugs, count, err
		}
		cellDrug.Datasets = strings.Split(datasets, ",")
		cellDrugs = append(cellDrugs, cellDrug)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogPrivateError(err)
		return cellDrugs, count, err
	}
	return cellDrugs, count, nil
}
