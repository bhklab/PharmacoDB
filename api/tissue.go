package api

import (
	"database/sql"
	"fmt"
	"strings"
)

// TD models experiment data for
// a tissue and a drug.
type TD struct {
	Drug     Drug     `json:"drug"`
	Datasets []string `json:"datasets"`
	Count    int      `json:"experiment_count"`
}

// TDS is a collection of TD.
type TDS []TD

// List updates receiver with a list of all tissues without pagination.
func (tissues *Tissues) List() error {
	var tissue Tissue
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT tissue_id, tissue_name FROM tissues;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*tissues = append(*tissues, tissue)
	}
	return nil
}

// ListPaginated updates receiver with a list of tissues using pagination.
func (tissues *Tissues) ListPaginated(page int, limit int) error {
	var tissue Tissue
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT tissue_id, tissue_name FROM tissues LIMIT %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*tissues = append(*tissues, tissue)
	}
	return nil
}

// Find updates receiver with a record for a single tissue.
func (tissue *Tissue) Find(id string, typ string) error {
	var query string
	db, err := InitDB()
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
			LogPrivateError(err)
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
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	query := "SELECT s.tissue_name, d.dataset_name FROM source_tissue_names s JOIN datasets d ON d.dataset_id = s.source_id WHERE s.tissue_id = ?;"
	rows, err := db.Query(query, tissue.ID)
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
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return cells, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS cell_id, cell_name FROM cells WHERE tissue_id = ? LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, tissue.ID)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return cells, count, err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Name)
		if err != nil {
			LogPrivateError(err)
			return cells, count, err
		}
		cells = append(cells, cell)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogPrivateError(err)
		return cells, count, err
	}
	return cells, count, nil
}

// Drugs returns a paginated list of all distinct drugs which have been tested with tissue, along with
// experiments count and an array of datasets that tested each cell/drug combination.
func (tissue *Tissue) Drugs(page int, limit int) (TDS, int, error) {
	var (
		tissueDrug  TD
		tissueDrugs TDS
		datasets    string
		count       int
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return tissueDrugs, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS d.drug_id, d.drug_name, GROUP_CONCAT(DISTINCT da.dataset_name) AS datasets, COUNT(*) AS experiment_count FROM experiments e JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.tissue_id = ? GROUP BY e.drug_id ORDER BY experiment_count DESC LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, tissue.ID)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return tissueDrugs, count, err
	}
	for rows.Next() {
		err = rows.Scan(&tissueDrug.Drug.ID, &tissueDrug.Drug.Name, &datasets, &tissueDrug.Count)
		if err != nil {
			LogPrivateError(err)
			return tissueDrugs, count, err
		}
		tissueDrug.Datasets = strings.Split(datasets, ",")
		tissueDrugs = append(tissueDrugs, tissueDrug)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogPrivateError(err)
		return tissueDrugs, count, err
	}
	return tissueDrugs, count, nil
}
