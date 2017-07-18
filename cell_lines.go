package main

import "fmt"

// ListPaginatedCells returns a list of paginated cell lines.
func ListPaginatedCells(page int, limit int) (Cells, error) {
	var (
		cell  Cell
		cells Cells
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT cell_id, accession_id, cell_name FROM cells LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return nil, err
		}
		cells = append(cells, cell)
	}
	return cells, nil
}

// ListAllCells returns a list of all cell lines without pagination.
func ListAllCells() (Cells, error) {
	var (
		cell  Cell
		cells Cells
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT cell_id, accession_id, cell_name FROM cells;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return nil, err
		}
		cells = append(cells, cell)
	}
	return cells, nil
}
