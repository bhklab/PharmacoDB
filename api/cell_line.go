package api

// Cell is a model for a cell line.
type Cell struct {
	ID   int     `json:"id"`
	ACC  *string `json:"accession_id,omitempty"`
	Name string  `json:"name"`
	// Tissue      *Tissue     `json:"tissue,omitempty"`
	// Annotations Annotations `json:"annotations,omitempty"`
}

// Cells is a collection of Cell.
type Cells []Cell

// List updates cells with a list of all cell lines without pagination.
func (cells *Cells) List() error {
	var cell Cell
	db, err := InitDB()
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
