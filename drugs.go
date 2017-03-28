package main

import (
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetDrugs handles GET requests for /drugs endpoint.
func GetDrugs(c *gin.Context) {
	getDataTypes(c, "List of all drugs in pharmacodb", "select drug_id, drug_name from drugs;")
}

// GetDrugStats handles GET requests for /drugs/stats endpoint.
func GetDrugStats(c *gin.Context) {
	getDataTypeStats(c, "Number of drugs tested in each dataset", "select dataset_id, drugs from dataset_statistics;")
}

// GetDrugIDs handles GET requests for /drugs/ids endpoint.
func GetDrugIDs(c *gin.Context) {
	getDataTypeIDs(c, "List of all drug IDs in pharmacodb", "select drug_id from drugs;")
}

// GetDrugNames handles GET requests for /drugs/names endpoint.
func GetDrugNames(c *gin.Context) {
	getDataTypeNames(c, "List of all drug names in pharmacodb", "select drug_name from drugs;")
}

// getDrug finds a drug using either ID or name.
func getDrug(c *gin.Context, ptype string) {
	var (
		drug      Drug
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
		queryStr = "select d.drug_id, d.drug_name, s.source_name, sdn.drug_name from drugs d inner join source_drug_names sdn on sdn.drug_id = d.drug_id inner join sources s on s.source_id = sdn.source_id where d.drug_id = ?"
	} else {
		queryStr = "select d.drug_id, d.drug_name, s.source_name, sdn.drug_name from drugs d inner join source_drug_names sdn on sdn.drug_id = d.drug_id inner join sources s on s.source_id = sdn.source_id where d.drug_name = ?"
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
		err = rows.Scan(&drug.ID, &drug.Name, &synsource, &syname)
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
		handleError(c, nil, http.StatusNotFound, fmt.Sprintf("Tissue with %s - %s - not found in pharmacodb", ptype, iden))
		return
	}
	drug.Synonyms = syns

	c.IndentedJSON(http.StatusOK, gin.H{
		"type": "drug",
		"data": drug,
	})
}

// GetDrugByID handles GET requests for /drugs/ids/:id endpoint.
func GetDrugByID(c *gin.Context) {
	getDrug(c, "id")
}

// GetDrugByName handles GET requests for /drugs/names/:name endpoint.
func GetDrugByName(c *gin.Context) {
	getDrug(c, "name")
}

// getDrugCells finds cell lines tested with a drug across datasets.
func getDrugCells(c *gin.Context, ptype string) {
	var (
		queryStr string
		cell     string
		dataset  string
		dcells   []CellDataset
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	iden := c.Param(ptype)
	if ptype == "id" {
		queryStr = "select c.cell_name, da.dataset_name from drugs dr inner join experiments e on e.drug_id = dr.drug_id inner join cells c on c.cell_id = e.cell_id inner join datasets da on da.dataset_id = e.dataset_id where dr.drug_id = ?"
	} else {
		queryStr = "select c.cell_name, da.dataset_name from drugs dr inner join experiments e on e.drug_id = dr.drug_id inner join cells c on c.cell_id = e.cell_id inner join datasets da on da.dataset_id = e.dataset_id where dr.drug_name = ?"
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
		err = rows.Scan(&cell, &dataset)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[cell] {
			for i, dc := range dcells {
				if dc.Cell == cell {
					if arrMember(dataset, dc.Datasets) {
						dcells[i].Experiments++
						break
					}
					dcells[i].Datasets = append(dcells[i].Datasets, dataset)
					dcells[i].Experiments++
					break
				}
			}
		} else {
			var dcell CellDataset
			dcell.Cell = cell
			dcell.Datasets = append(dcell.Datasets, dataset)
			dcell.Experiments = 1
			dcells = append(dcells, dcell)
			exists[cell] = true
		}
		iter = 1
	}
	if iter == 0 {
		handleError(c, nil, http.StatusNotFound, "No cell lines found tested with this drug")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": "List of cell lines tested with drug across datasets",
		"count":       len(dcells),
		"data":        dcells,
	})
}

// GetDrugCellsByID handles GET requests for /drugs/ids/:id/cell_lines endpoint.
func GetDrugCellsByID(c *gin.Context) {
	getDrugCells(c, "id")
}

// GetDrugCellsByName handles GET requests for /drugs/names/:name/cell_lines endpoint.
func GetDrugCellsByName(c *gin.Context) {
	getDrugCells(c, "name")
}

// getDrugTissues finds tissues tested with a drug across datasets.
func getDrugTissues(c *gin.Context, ptype string) {
	var (
		queryStr string
		tissue   string
		dataset  string
		dtissues []TissueDataset
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	iden := c.Param(ptype)
	if ptype == "id" {
		queryStr = "select t.tissue_name, da.dataset_name from drugs dr inner join experiments e on e.drug_id = dr.drug_id inner join tissues t on t.tissue_id = e.tissue_id inner join datasets da on da.dataset_id = e.dataset_id where dr.drug_id = ?"
	} else {
		queryStr = "select t.tissue_name, da.dataset_name from drugs dr inner join experiments e on e.drug_id = dr.drug_id inner join tissues t on t.tissue_id = e.tissue_id inner join datasets da on da.dataset_id = e.dataset_id where dr.drug_name = ?"
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
		err = rows.Scan(&tissue, &dataset)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if exists[tissue] {
			for i, dt := range dtissues {
				if dt.Tissue == tissue {
					if arrMember(dataset, dt.Datasets) {
						dtissues[i].Experiments++
						break
					}
					dtissues[i].Datasets = append(dtissues[i].Datasets, dataset)
					dtissues[i].Experiments++
					break
				}
			}
		} else {
			var dtissue TissueDataset
			dtissue.Tissue = tissue
			dtissue.Datasets = append(dtissue.Datasets, dataset)
			dtissue.Experiments = 1
			dtissues = append(dtissues, dtissue)
			exists[tissue] = true
		}
		iter = 1
	}
	if iter == 0 {
		handleError(c, nil, http.StatusNotFound, "No tissue found tested with this drug")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"description": "List of tissues tested with drug across datasets",
		"count":       len(dtissues),
		"data":        dtissues,
	})
}

// GetDrugTissuesByID handles GET requests for /drugs/ids/:id/tissues endpoint.
func GetDrugTissuesByID(c *gin.Context) {
	getDrugTissues(c, "id")
}

// GetDrugTissuesByName handles GET requests for /drugs/names/:name/tissues endpoint.
func GetDrugTissuesByName(c *gin.Context) {
	getDrugTissues(c, "name")
}
