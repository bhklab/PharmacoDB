package api

import (
	"database/sql"
	"fmt"
)

// List updates receiver with a list of all datasets without pagination.
func (datasets *Datasets) List() error {
	var dataset Dataset
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT dataset_id, dataset_name FROM datasets;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&dataset.ID, &dataset.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*datasets = append(*datasets, dataset)
	}
	return nil
}

// ListPaginated updates receiver with a list of datasets using pagination.
func (datasets *Datasets) ListPaginated(page int, limit int) error {
	var dataset Dataset
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT dataset_id, dataset_name FROM datasets LIMIT %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&dataset.ID, &dataset.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*datasets = append(*datasets, dataset)
	}
	return nil
}

// Find updates receiver with a record for a single dataset.
func (dataset *Dataset) Find(id string, typ string) error {
	var query string
	db, err := InitDB()
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
			LogPrivateError(err)
		}
		return err
	}
	return nil
}
