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
				"message":           "Bad request. Check PharmacoDb API official documentation for properly formed endpoints/routes.",
				"documentation_url": "https://www.pharmacodb.com/docs/api",
			},
		})
	})

	v1 := router.Group("v1")
	{
		v1.GET("/datatypes", GetDataTypes)
		v1.GET("/cell_lines", GetCells)
		v1.GET("/tissues", GetTissues)
		v1.GET("/drugs", GetDrugs)
		v1.GET("/datasets", GetDatasets)
	}

	router.Run(":3000")
}
