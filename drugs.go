package main

import "gopkg.in/gin-gonic/gin.v1"

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
