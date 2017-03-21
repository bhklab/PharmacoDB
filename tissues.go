package main

import "gopkg.in/gin-gonic/gin.v1"

// GetTissues handles GET requests for /tissues endpoint.
func GetTissues(c *gin.Context) {
	queryStr := "select tissue_id, tissue_name from tissues;"
	desc := "List of all tissues in pharmacodb"
	getDataTypes(c, desc, queryStr)
}

// GetTissueStats handles GET requests for /tissues/stats endpoint.
func GetTissueStats(c *gin.Context) {
	queryStr := "select dataset_id, tissues from dataset_statistics;"
	desc := "Number of tissues tested in each dataset"
	getDataTypeStats(c, desc, queryStr)
}
