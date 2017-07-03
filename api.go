package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.StaticFile("/favicon.ico", "./lib/foobar.png")

	router.GET("/", func(c *gin.Context) {
		welcome := "Welcome to PharmacoDB API."
		link := "For more information, visit: https://www.pharmacodb.com/docs/api"
		message := welcome + "\n" + link
		c.String(http.StatusOK, message)
	})

	v1 := router.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			info := "To use the API, add endpoints to the request."
			link := "For more information, visit: https://www.pharmacodb.com/docs/api"
			message := info + "\n" + link
			c.String(http.StatusOK, message)
		})

		v1.GET("/cell_lines", IndexCell)

		v1.GET("/tissues", IndexCell)

		v1.GET("/drugs", IndexCell)

		v1.GET("/datasets", IndexCell)
	}

	// Responds with status code 400 (Bad Request) if no routers match the request url.
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":              http.StatusBadRequest,
				"message":           "Bad Request",
				"documentation_url": "https://www.pharmacodb.com/docs/api",
			},
		})
	})

	// Listen and serve on 0.0.0.0:4200
	router.Run(":4200")
}
