package api

import (
	"database/sql"
	"fmt"
)

// List updates a list of tissues without pagination.
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

// ListPaginated returns a list of paginated tissues.
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

// Find finds and updates tissue with a record for a Tissue.
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

// Annotate adds annotations to a tissue.
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
