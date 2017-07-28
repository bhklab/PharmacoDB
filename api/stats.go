package api

// TissueCount models the number of cell lines per tissue.
type TissueCount struct {
	Tissue Tissue `json:"tissue"`
	Count  int    `json:"cell_lines_count"`
}

// TissueCounts is a collection of TissueCount.
type TissueCounts []TissueCount

// DatasetCount models the number of drugs tested per dataset.
type DatasetCount struct {
	Dataset Dataset `json:"dataset"`
	Count   int     `json:"drugs_count"`
}

// DatasetCounts is a collection of DatasetDrugCount.
type DatasetCounts []DatasetCount

// CountCellsPerTissue returns a list of all tissues, along with the number of
// cell lines of each tissue type.
func CountCellsPerTissue() (TissueCounts, error) {
	var (
		tissueCellCount  TissueCount
		tissueCellCounts TissueCounts
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return tissueCellCounts, err
	}
	query := "SELECT t.tissue_id, t.tissue_name, COUNT(*) AS cell_lines_count FROM tissues t JOIN cells c ON c.tissue_id = t.tissue_id GROUP BY(c.tissue_id);"
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return tissueCellCounts, err
	}
	for rows.Next() {
		err = rows.Scan(&tissueCellCount.Tissue.ID, &tissueCellCount.Tissue.Name, &tissueCellCount.Count)
		if err != nil {
			LogPrivateError(err)
			return tissueCellCounts, err
		}
		tissueCellCounts = append(tissueCellCounts, tissueCellCount)
	}
	return tissueCellCounts, nil
}

// CountItemsPerDataset returns a list of all datasets, along with the number of
// required item tested in each dataset.
func CountItemsPerDataset(query string) (DatasetCounts, error) {
	var (
		count  DatasetCount
		counts DatasetCounts
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return counts, err
	}
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return counts, err
	}
	for rows.Next() {
		err = rows.Scan(&count.Dataset.ID, &count.Dataset.Name, &count.Count)
		if err != nil {
			LogPrivateError(err)
			return counts, err
		}
		counts = append(counts, count)
	}
	return counts, nil
}
