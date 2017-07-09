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
		v1.GET("/cell_lines/:id", ShowCell)

		v1.GET("/tissues", IndexTissue)
		// v1.GET("/tissues/:id", ShowTissue)

		v1.GET("/drugs", IndexDrug)
		// v1.GET("/drugs/:id", ShowDrug)

		v1.GET("/datasets", IndexDataset)
		// v1.GET("/datasets/:id", ShowDataset)

		v1.GET("/experiments", IndexExperiment)
		v1.GET("/experiments/:id", ShowExperiment)
	}

	// Responds with status code 400 (Bad Request) if no routers match the request url.
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
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
