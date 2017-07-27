package api

import "strconv"

// TissueCellCount models the number of cell lines per tissue.
type TissueCellCount struct {
	Tissue Tissue `json:"tissue"`
	Count  int    `json:"cell_lines_count"`
}

// TissueCellCounts is a collection of TissueCellsCount.
type TissueCellCounts []TissueCellCount

// DatasetDrugCount models the number of drugs tested per dataset.
type DatasetDrugCount struct {
	Dataset Dataset `json:"dataset"`
	Count   int     `json:"drugs_count"`
}

// DatasetDrugCounts is a collection of DatasetDrugCount.
type DatasetDrugCounts []DatasetDrugCount

// CountCellsPerTissue returns a list of all tissues, along with the number of
// cell lines of each tissue type.
func CountCellsPerTissue() (TissueCellCounts, error) {
	var (
		tissueCellCount  TissueCellCount
		tissueCellCounts TissueCellCounts
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

// CountDrugsPerDataset returns a list of all datasets, along with the number of
// drugs tested in each dataset.
func CountDrugsPerDataset() (DatasetDrugCounts, error) {
	type DDD struct {
		ID    int
		Count int
	}
	var (
		datasetDrugCount  DatasetDrugCount
		datasetDrugCounts DatasetDrugCounts
		DD                DDD
		DDs               []DDD
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return datasetDrugCounts, err
	}
	query := "SELECT dataset_id, COUNT(DISTINCT drug_id) AS drugs_count FROM experiments GROUP BY dataset_id;"
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(err)
		return datasetDrugCounts, err
	}
	for rows.Next() {
		err = rows.Scan(&DD.ID, &DD.Count)
		if err != nil {
			LogPrivateError(err)
			return datasetDrugCounts, err
		}
		DDs = append(DDs, DD)
	}
	for _, a := range DDs {
		var dataset Dataset
		err = dataset.Find(strconv.Itoa(a.ID), "id")
		if err != nil {
			return datasetDrugCounts, err
		}
		datasetDrugCount.Dataset = dataset
		datasetDrugCount.Count = a.Count
		datasetDrugCounts = append(datasetDrugCounts, datasetDrugCount)
	}
	return datasetDrugCounts, nil
}
