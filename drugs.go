package main

import (
	"database/sql"
	"fmt"
)

// PaginatedDrugs returns a list of paginated drugs.
func PaginatedDrugs(page int, limit int) (Drugs, error) {
	var (
		drug  Drug
		drugs Drugs
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return drugs, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT drug_id, drug_name FROM drugs LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return drugs, err
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return drugs, err
		}
		drugs = append(drugs, drug)
	}
	return drugs, nil
}

// NonPaginatedDrugs returns a list of all drugs without pagination.
func NonPaginatedDrugs() (Drugs, error) {
	var (
		drug  Drug
		drugs Drugs
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return drugs, err
	}
	rows, err := db.Query("SELECT drug_id, drug_name FROM drugs;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return drugs, err
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return drugs, err
		}
		drugs = append(drugs, drug)
	}
	return drugs, nil
}

// FindDrug returns a drug, queried using id or name.
func FindDrug(id string, typ string) (Drug, error) {
	var (
		drug  Drug
		query string
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return drug, err
	}
	if sameString(typ, "name") {
		query = "SELECT drug_id, drug_name FROM drugs WHERE drug_name LIKE ?;"
	} else {
		query = "SELECT drug_id, drug_name FROM drugs WHERE drug_id LIKE ?;"
	}
	row := db.QueryRow(query, id)
	err = row.Scan(&drug.ID, &drug.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogPrivateError(ErrorTypePrivate, err)
		}
		return drug, err
	}
	return drug, nil
}

// Annotate adds annotations to a drug.
func (drug *Drug) Annotate() error {
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
	query := "SELECT s.drug_name, d.dataset_name FROM source_drug_names s JOIN datasets d ON d.dataset_id = s.source_id WHERE s.drug_id = ?;"
	rows, err := db.Query(query, drug.ID)
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
	drug.Annotations = annotations
	return nil
}
