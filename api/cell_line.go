package api

import (
	"database/sql"
	"fmt"
)

// List updates cells with a list of all cell lines without pagination.
func (cells *Cells) List() error {
	var cell Cell
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT cell_id, accession_id, cell_name FROM cells;")
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*cells = append(*cells, cell)
	}
	return nil
}

// ListPaginated returns a list of paginated cell lines.
func (cells *Cells) ListPaginated(page int, limit int) error {
	var cell Cell
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	s := (page - 1) * limit
	rows, err := db.Query(fmt.Sprintf("SELECT cell_id, accession_id, cell_name FROM cells LIMIT %d,%d;", s, limit))
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.ACC, &cell.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*cells = append(*cells, cell)
	}
	return nil
}

// Find finds and updates cell with a record for a cell line.
func (cell *Cell) Find(id string, typ string) error {
	var query string
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	if isSameStr(typ, "name") {
		query = "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE c.cell_name LIKE ?;"
	} else if isSameStr(typ, "accession") {
		query = "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE c.accession_id LIKE ?;"
	} else {
		query = "SELECT c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name FROM cells c JOIN tissues t ON t.tissue_id = c.tissue_id WHERE c.cell_id LIKE ?;"
	}
	tissue := &Tissue{}
	row := db.QueryRow(query, id)
	err = row.Scan(&cell.ID, &cell.ACC, &cell.Name, &tissue.ID, &tissue.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogPrivateError(err)
		}
		return err
	}
	cell.Tissue = tissue
	return nil
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
	cell.Annotations = annotations
	return nil
}
