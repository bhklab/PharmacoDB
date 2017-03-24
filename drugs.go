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
