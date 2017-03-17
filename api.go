package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(400, gin.H{
			"error": gin.H{
				"status":            400,
				"message":           "Bad Request",
				"suggestions":       "Check the official API documentation to see how to properly format endpoints/routes",
				"documentation_url": "https://www.pharmacodb.com/docs/api",
			},
		})
	})

	v1 := router.Group("v1")
	{
		v1.GET("/cell_lines", GetCells)
		v1.GET("/cell_lines/stats", GetCellStats)

		v1.GET("/tissues", GetTissues)
		v1.GET("/tissues/stats", GetTissueStats)

		v1.GET("/drugs", GetDrugs)
		v1.GET("/drugs/stats", GetDrugStats)

		v1.GET("/datasets", GetDatasets)
		v1.GET("/datasets/stats", GetDatasetStats)
	}

	router.Run(":3000")
}
