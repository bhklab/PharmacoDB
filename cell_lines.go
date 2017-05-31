package main

import (
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetCells handles GET requests for /cell_lines endpoint.
func GetCells(c *gin.Context) {
	getDataTypes(c, "List of all cell lines in pharmacodb", "select cell_id, cell_name from cells;")
}

// GetCellStats handles GET requests for /cell_lines/stats endpoint.
func GetCellStats(c *gin.Context) {
	getDataTypeStats(c, "Number of cell lines tested in each dataset", "select dataset_id, cell_lines from dataset_statistics;")
}

// GetCellIDs handles GET requests for /cell_lines/ids endpoint.
func GetCellIDs(c *gin.Context) {
	getDataTypeIDs(c, "List of all cell line IDs in pharmacodb", "select cell_id from cells;")
}

// GetCellNames handles GET requests for /cell_lines/names endpoint.
func GetCellNames(c *gin.Context) {
	getDataTypeNames(c, "List of all cell line names in pharmacodb", "select cell_name from cells;")
}

// getCell finds a cell line using either ID or name.
func getCell(c *gin.Context, ptype string) {
	var (
		cell      Cell
		syname    string
		synsource string
		syns      []Synonym
		queryStr  string
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	iden := c.Param(ptype)
	if ptype == "id" {
		queryStr = "select c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name, s.source_name, scn.cell_name from cells c inner join tissues t on t.tissue_id = c.tissue_id inner join source_cell_names scn on scn.cell_id = c.cell_id inner join sources s on s.source_id = scn.source_id where c.cell_id = ?"
	} else {
		queryStr = "select c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name, s.source_name, scn.cell_name from cells c inner join tissues t on t.tissue_id = c.tissue_id inner join source_cell_names scn on scn.cell_id = c.cell_id inner join sources s on s.source_id = scn.source_id where c.cell_name = ?"
	}
	rows, err := db.Query(queryStr, iden)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[string]bool)
	iter := 0
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Accession, &cell.Name, &cell.Tissue.ID, &cell.Tissue.Name, &synsource, &syname)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[syname] {
			for i, syn := range syns {
				if syn.Name == syname {
					syns[i].Datasets = append(syns[i].Datasets, synsource)
					break
				}
			}
		} else {
			var syn Synonym
			syn.Name = syname
			syn.Datasets = append(syn.Datasets, synsource)
			syns = append(syns, syn)
			exists[syname] = true
		}
		iter = 1
	}
	if iter == 0 {
		handleError(c, nil, http.StatusNotFound, fmt.Sprintf("Cell line with %s - %s - not found in pharmacodb", ptype, iden))
		return
	}
	cell.Synonyms = syns

	c.IndentedJSON(http.StatusOK, gin.H{
		"type": "cell line",
		"data": cell,
	})
}

// GetCellByID handles GET requests for /cell_lines/ids/:id endpoint.
func GetCellByID(c *gin.Context) {
	getCell(c, "id")
}

// GetCellByName handles GET requests for /cell_lines/names/:name endpoint.
func GetCellByName(c *gin.Context) {
	getCell(c, "name")
}

// getCellDrugs finds all drugs tested with a cell line across datasets.
func getCellDrugs(c *gin.Context, ptype string) {
	var (
		queryStr string
		drug     string
		dataset  string
		cdrugs   []DrugDataset
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	iden := c.Param(ptype)
	if ptype == "id" {
		queryStr = "select dr.drug_name, da.dataset_name from experiments e inner join drugs dr on dr.drug_id = e.drug_id inner join datasets da on da.dataset_id = e.dataset_id where e.cell_id = ?"
	} else {
		queryStr = "select dr.drug_id, da.dataset_id from cells c inner join experiments e on e.cell_id = c.cell_id inner join drugs dr on dr.drug_id = e.drug_id inner join datasets da on da.dataset_id = e.dataset_id where c.cell_name = ?"
	}
	rows, err := db.Query(queryStr, iden)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[string]bool)
	iter := 0
	for rows.Next() {
		err = rows.Scan(&drug, &dataset)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[drug] {
			for i, cd := range cdrugs {
				if cd.Drug == drug {
					if arrMember(dataset, cd.Datasets) {
						cdrugs[i].Experiments++
						break
					}
					cdrugs[i].Datasets = append(cdrugs[i].Datasets, dataset)
					cdrugs[i].Experiments++
					break
				}
			}
		} else {
			var cdrug DrugDataset
			cdrug.Drug = drug
			cdrug.Datasets = append(cdrug.Datasets, dataset)
			cdrug.Experiments = 1
			cdrugs = append(cdrugs, cdrug)
			exists[drug] = true
		}
		iter = 1
	}
	if iter == 0 {
		handleError(c, nil, http.StatusNotFound, "No drugs found tested with this cell line")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": "List of drugs tested with cell line across datasets",
		"count":       len(cdrugs),
		"data":        cdrugs,
	})
}

// GetCellDrugsByID handles GET requests for /cell_lines/ids/:id/drugs endpoint.
func GetCellDrugsByID(c *gin.Context) {
	getCellDrugs(c, "id")
}

// GetCellDrugsByName handles GET requests for /cell_lines/names/:name/drugs endpoint.
func GetCellDrugsByName(c *gin.Context) {
	getCellDrugs(c, "name")
}

// getCellDrugStats finds number of drugs tested with cell line in datasets.
// (currently, datasets with zero experiments are not returned by this method)
func getCellDrugStats(c *gin.Context, ptype string) {
	var (
		drug     string
		dataset  string
		queryStr string
		dstats   []DatasetStat
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	iden := c.Param(ptype)
	if ptype == "id" {
		queryStr = "select distinct drug_id, d.dataset_name from experiments e inner join datasets d on d.dataset_id = e.dataset_id where cell_id = ?;"
	} else {
		queryStr = "select distinct e.drug_id, d.dataset_name from cells c inner join experiments e on e.cell_id = c.cell_id inner join datasets d on d.dataset_id = e.dataset_id where c.cell_name = ?;"
	}
	rows, err := db.Query(queryStr, iden)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[string]bool)
	iter := 0
	for rows.Next() {
		err = rows.Scan(&drug, &dataset)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[dataset] {
			for i, d := range dstats {
				if d.Dataset == dataset {
					dstats[i].Count++
					break
				}
			}
		} else {
			var dstat DatasetStat
			dstat.Dataset = dataset
			dstat.Count = 1
			dstats = append(dstats, dstat)
			exists[dataset] = true
		}
		iter = 1
	}
	if iter == 0 {
		handleError(c, nil, http.StatusNotFound, "No drugs found tested with this cell line")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": "Number of drugs tested with cell line in datasets",
		"data":        dstats,
	})
}

// GetCellDrugStatsByID handles GET requests for /cell_lines/ids/:id/drugs/stats endpoint.
func GetCellDrugStatsByID(c *gin.Context) {
	getCellDrugStats(c, "id")
}

// GetCellDrugStatsByName handles GET requests for /cell_lines/names/:name/drugs/stats endpoint.
func GetCellDrugStatsByName(c *gin.Context) {
	getCellDrugStats(c, "name")
}
