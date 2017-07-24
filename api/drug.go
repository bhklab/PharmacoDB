package api

import (
	"database/sql"
	"fmt"
	"strings"
)

// DC models experiment data for
// a drug and a cell line.
type DC struct {
	Cell     Cell     `json:"cell_line"`
	Datasets []string `json:"datasets"`
	Count    int      `json:"experiment_count"`
}

// DCS is a collection of DCT.
type DCS []DC

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

// Cells returns a paginated list of all distinct cell lines which have been tested with drug, along with
// experiments count and an array of datasets that tested each cell/drug combination.
func (drug *Drug) Cells(page int, limit int) (DCS, int, error) {
	var (
		drugCell  DC
		drugCells DCS
		datasets  string
		count     int
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return drugCells, count, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT SQL_CALC_FOUND_ROWS c.cell_id, c.cell_name, GROUP_CONCAT(DISTINCT da.dataset_name) AS datasets, COUNT(*) AS experiment_count FROM experiments e JOIN cells c ON c.cell_id = e.cell_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.drug_id = ? GROUP BY e.cell_id ORDER BY experiment_count DESC LIMIT %d,%d;", s, limit)
	rows, err := db.Query(query, drug.ID)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return drugCells, count, err
	}
	for rows.Next() {
		err = rows.Scan(&drugCell.Cell.ID, &drugCell.Cell.Name, &datasets, &drugCell.Count)
		if err != nil {
			LogPrivateError(err)
			return drugCells, count, err
		}
		drugCell.Datasets = strings.Split(datasets, ",")
		drugCells = append(drugCells, drugCell)
	}
	row := db.QueryRow("SELECT FOUND_ROWS();")
	err = row.Scan(&count)
	if err != nil {
		LogPrivateError(err)
		return drugCells, count, err
	}
	return drugCells, count, nil
}
