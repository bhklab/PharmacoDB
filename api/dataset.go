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

// Cells returns a paginated list of all cell lines tested in a dataset.
func (dataset *Dataset) Cells(page int, limit int) (Cells, int, error) {
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
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS c.cell_id, c.cell_name FROM experiments e JOIN cells c ON c.cell_id = e.cell_id WHERE e.dataset_id = ? GROUP BY(e.cell_id) LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, dataset.ID)
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

// Tissues returns a paginated list of all tissues tested in dataset.
func (dataset *Dataset) Tissues(page int, limit int) (Tissues, int, error) {
	var (
		tissue  Tissue
		tissues Tissues
		count   int
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return tissues, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS t.tissue_id, t.tissue_name FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id WHERE e.dataset_id = ? GROUP BY(e.tissue_id) LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, dataset.ID)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return tissues, count, err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogPrivateError(err)
			return tissues, count, err
		}
		tissues = append(tissues, tissue)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogPrivateError(err)
		return tissues, count, err
	}
	return tissues, count, nil
}

// Drugs returns a paginated list of all drugs tested in dataset.
func (dataset *Dataset) Drugs(page int, limit int) (Drugs, int, error) {
	var (
		drug  Drug
		drugs Drugs
		count int
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return drugs, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS d.drug_id, d.drug_name FROM experiments e JOIN drugs d ON d.drug_id = e.drug_id WHERE e.dataset_id = ? GROUP BY(e.drug_id) LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, &dataset.ID)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return drugs, count, err
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, drug.Name)
		if err != nil {
			LogPrivateError(err)
			return drugs, count, err
		}
		drugs = append(drugs, drug)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogPrivateError(err)
		return drugs, count, err
	}
	return drugs, count, nil
}
