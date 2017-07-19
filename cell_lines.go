package main

import (
	"database/sql"
	"fmt"
)

// PaginatedCells returns a list of paginated cell lines.
func PaginatedCells(page int, limit int) (Cells, error) {
	var (
		cell  Cell
		cells Cells
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return cells, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT cell_id, accession_id, cell_name FROM cells LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return cells, err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return cells, err
		}
		cells = append(cells, cell)
	}
	return cells, nil
}

// NonPaginatedCells returns a list of all cell lines without pagination.
func NonPaginatedCells() (Cells, error) {
	var (
		cell  Cell
		cells Cells
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return cells, err
	}
	rows, err := db.Query("SELECT cell_id, accession_id, cell_name FROM cells;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return cells, err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return cells, err
		}
		cells = append(cells, cell)
	}
	return cells, nil
}

// FindCell returns a cell line, queried using id, name or accession.
func FindCell(id string, typ string) (Cell, error) {
	var (
		cell  Cell
		query string
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return cell, err
	}
	if sameString(typ, "name") {
		query = "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE c.cell_name LIKE ?;"
	} else if sameString(typ, "accession") {
		query = "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE c.accession_id LIKE ?;"
	} else {
		query = "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE c.cell_id LIKE ?;"
	}
	tissue := &Tissue{}
	row := db.QueryRow(query, id)
	err = row.Scan(&cell.ID, &cell.ACC, &cell.Name, &tissue.ID, &tissue.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogPrivateError(ErrorTypePrivate, err)
		}
		return cell, err
	}
	cell.Tissue = tissue
	return cell, nil
}

// Annotate adds annotations to a cell line.
func (cell *Cell) Annotate() error {
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
	query := "SELECT s.cell_name, d.dataset_name FROM source_cell_names s JOIN datasets d ON d.dataset_id = s.source_id WHERE s.cell_id = ?;"
	rows, err := db.Query(query, cell.ID)
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
	cell.Annotations = annotations
	return nil
}
