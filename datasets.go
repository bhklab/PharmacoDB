package main

import "fmt"

// ListPaginatedDatasets returns a list of paginated datasets.
func ListPaginatedDatasets(page int, limit int) (Datasets, error) {
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

// ListAllDatasets returns a list of all datasets without pagination.
func ListAllDatasets() (Datasets, error) {
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
