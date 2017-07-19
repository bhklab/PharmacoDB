package main

import (
	"database/sql"
	"fmt"
)

// PaginatedDatasets returns a list of paginated datasets.
func PaginatedDatasets(page int, limit int) (Datasets, error) {
	var (
		dataset  Dataset
		datasets Datasets
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT dataset_id, dataset_name FROM datasets LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&dataset.ID, &dataset.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return nil, err
		}
		datasets = append(datasets, dataset)
	}
	return datasets, nil
}

// NonPaginatedDatasets returns a list of all datasets without pagination.
func NonPaginatedDatasets() (Datasets, error) {
	var (
		dataset  Dataset
		datasets Datasets
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT dataset_id, dataset_name FROM datasets")
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&dataset.ID, &dataset.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return nil, err
		}
		datasets = append(datasets, dataset)
	}
	return datasets, nil
}

// FindDataset returns a dataset, queried using id or name.
func FindDataset(id string, typ string) (Dataset, error) {
	var (
		dataset Dataset
		query   string
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return dataset, err
	}
	if sameString(typ, "name") {
		query = "SELECT dataset_id, dataset_name FROM datasets WHERE dataset_name LIKE ?;"
	} else {
		query = "SELECT dataset_id, dataset_name FROM datasets WHERE dataset_id LIKE ?;"
	}
	row := db.QueryRow(query, id)
	err = row.Scan(&dataset.ID, &dataset.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogPrivateError(ErrorTypePrivate, err)
		}
		return dataset, err
	}
	return dataset, nil
}
