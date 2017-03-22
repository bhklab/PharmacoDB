package main

import "gopkg.in/gin-gonic/gin.v1"

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
