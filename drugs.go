package main

import "fmt"

// ListPaginatedDrugs returns a list of paginated drugs.
func ListPaginatedDrugs(page int, limit int) (Drugs, error) {
	var (
		drug  Drug
		drugs Drugs
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT drug_id, drug_name FROM drugs LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return nil, err
		}
		drugs = append(drugs, drug)
	}
	return drugs, nil
}

// ListAllDrugs returns a list of all drugs without pagination.
func ListAllDrugs() (Drugs, error) {
	var (
		drug  Drug
		drugs Drugs
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT drug_id, drug_name FROM drugs;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&drug.ID, &drug.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return nil, err
		}
		drugs = append(drugs, drug)
	}
	return drugs, nil
}
