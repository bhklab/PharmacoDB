package api

import (
	"database/sql"
	"fmt"
)

// List updates receiver with a list of all drugs without pagination.
func (drugs *Drugs) List() error {
	var drug Drug
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT drug_id, drug_name FROM drugs;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*drugs = append(*drugs, drug)
	}
	return nil
}

// ListPaginated updates receiver with a list of drugs using pagination.
func (drugs *Drugs) ListPaginated(page int, limit int) error {
	var drug Drug
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT drug_id, drug_name FROM drugs LIMIT %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*drugs = append(*drugs, drug)
	}
	return nil
}

// Find updates receiver with a record for a single Drug.
func (drug *Drug) Find(id string, typ string) error {
	var query string
	db, err := InitDB()
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
	err = row.Scan(&drug.ID, &drug.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogPrivateError(err)
		}
		return err
	}
	return nil
}

// Annotate adds annotations to drug.
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
	drug.Annotations = annotations
	return nil
}
