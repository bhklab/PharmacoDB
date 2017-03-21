package main

import "gopkg.in/gin-gonic/gin.v1"

// GetDrugs handles GET requests for /drugs endpoint.
func GetDrugs(c *gin.Context) {
	queryStr := "select drug_id, drug_name from drugs;"
	desc := "List of all drugs in pharmacodb"
	getDataTypes(c, desc, queryStr)
}

// GetDrugStats handles GET requests for /drugs/stats endpoint.
func GetDrugStats(c *gin.Context) {
	queryStr := "select dataset_id, drugs from dataset_statistics;"
	desc := "Number of drugs tested in each dataset"
	getDataTypeStats(c, desc, queryStr)
}
