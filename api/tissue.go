package api

import (
	"database/sql"
	"fmt"
)

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
