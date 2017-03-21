package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	router := gin.Default()

	router.StaticFile("/favicon.ico", "./favicon.ico")
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"status":            http.StatusBadRequest,
				"message":           "Bad Request",
				"suggestions":       "Check the official API documentation to see how to properly format endpoints/routes",
				"documentation_url": "https://www.pharmacodb.com/docs/api",
			},
		})
	})

	v1 := router.Group("v1")
	{
		v1.GET("/cell_lines", GetCells)
	}

	router.Run(":3000")
}
