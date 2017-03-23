package main

import (
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetTissues handles GET requests for /tissues endpoint.
func GetTissues(c *gin.Context) {
	getDataTypes(c, "List of all tissues in pharmacodb",
		"select tissue_id, tissue_name from tissues;")
}

// GetTissueStats handles GET requests for /tissues/stats endpoint.
func GetTissueStats(c *gin.Context) {
	getDataTypeStats(c, "Number of tissues tested in each dataset", "select dataset_id, tissues from dataset_statistics;")
}

// GetTissueIDs handles GET requests for /tissues/ids endpoint.
func GetTissueIDs(c *gin.Context) {
	getDataTypeIDs(c, "List of all tissue IDs in pharmacodb", "select tissue_id from tissues;")
}

// GetTissueNames handles GET requests for /tissues/names endpoint.
func GetTissueNames(c *gin.Context) {
	getDataTypeNames(c, "List of all tissue names in pharmacodb", "select tissue_name from tissues;")
}

// GetTissueByID handles GET requests for /tissues/ids/:id endpoint.
func GetTissueByID(c *gin.Context) {
	var (
		tissue    Tissue
		syname    string
		synsource string
		syns      []Synonym
	)

	db, err := initDB()
	defer db.Close()
	if err != nil {
		handleError(c, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	id := c.Param("id")
	queryStr := "select t.tissue_id, t.tissue_name, s.source_name, stn.tissue_name from tissues t inner join source_tissue_names stn on stn.tissue_id = t.tissue_id inner join sources s on s.source_id = stn.source_id where t.tissue_id = ?"
	rows, err := db.Query(queryStr, id)
	defer rows.Close()
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	exists := make(map[string]bool)
	iter := 0
	for rows.Next() {
		err = rows.Scan(&tissue.ID, &tissue.Name, &synsource, &syname)
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
		handleError(c, nil, http.StatusNotFound, fmt.Sprintf("Tissue with ID - %s - not found in pharmacodb", id))
		return
	}
	tissue.Synonyms = syns

	c.IndentedJSON(http.StatusOK, gin.H{
		"type": "tissue",
		"data": tissue,
	})
}
