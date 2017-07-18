package main

import "fmt"

// ListPaginatedTissues returns a list of paginated tissues.
func ListPaginatedTissues(page int, limit int) (Tissues, error) {
	var (
		tissue  Tissue
		tissues Tissues
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT tissue_id, tissue_name FROM tissues LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return nil, err
		}
		tissues = append(tissues, tissue)
	}
	return tissues, nil
}

// ListAllTissues returns a list of all tissues without pagination.
func ListAllTissues() (Tissues, error) {
	var (
		tissue  Tissue
		tissues Tissues
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT tissue_id, tissue_name FROM tissues;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return nil, err
		}
		tissues = append(tissues, tissue)
	}
	return tissues, nil
}
