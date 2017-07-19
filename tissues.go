package main

import (
	"database/sql"
	"fmt"
)

// PaginatedTissues returns a list of paginated tissues.
func PaginatedTissues(page int, limit int) (Tissues, error) {
	var (
		tissue  Tissue
		tissues Tissues
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return tissues, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT tissue_id, tissue_name FROM tissues LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return tissues, err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return tissues, err
		}
		tissues = append(tissues, tissue)
	}
	return tissues, nil
}

// NonPaginatedTissues returns a list of all tissues without pagination.
func NonPaginatedTissues() (Tissues, error) {
	var (
		tissue  Tissue
		tissues Tissues
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return tissues, err
	}
	rows, err := db.Query("SELECT tissue_id, tissue_name FROM tissues;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return tissues, err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return tissues, err
		}
		tissues = append(tissues, tissue)
	}
	return tissues, nil
}

// FindTissue returns a tissue, queried using ID or name.
func FindTissue(id string, typ string) (Tissue, error) {
	var (
		tissue Tissue
		query  string
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return tissue, err
	}
	if sameString(typ, "name") {
		query = "SELECT tissue_id, tissue_name FROM tissues WHERE tissue_name LIKE ?;"
	} else {
		query = "SELECT tissue_id, tissue_name FROM tissues WHERE tissue_id LIKE ?;"
	}
	row := db.QueryRow(query, id)
	err = row.Scan(&tissue.ID, &tissue.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogPrivateError(ErrorTypePrivate, err)
		}
		return tissue, err
	}
	return tissue, nil
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
		LogPrivateError(ErrorTypePrivate, err)
		return err
	}
	exists := make(map[string]bool)
	for rows.Next() {
		err = rows.Scan(&annotationName, &datasetName)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
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
