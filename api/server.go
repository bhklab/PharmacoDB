package api

import "github.com/gin-gonic/gin"

// Init server.
func Init(config *Config) {
	// Set database credentials
	SetDB(config.Version)

	SetVersion(config.Version)

	// Gin router with default middleware: logger and recovery
	router := gin.Default()

	// Set gin mode
	gin.SetMode(config.Mode)

	// Serve favicon
	router.StaticFile("/favicon.ico", "./assets/images/favicon.ico")

	router.GET("/", RootHandler)

	v := router.Group("v" + config.Version + "/")
	v.GET("/", RootHandler)
	// GET requests
	for _, route := range routesGET {
		v.Handle(GET, route.Endpoint, route.Handler)
	}
	// HEAD requests
	for _, route := range routesHEAD {
		v.Handle(HEAD, route.Endpoint, route.Handler)
	}

	// If no routers match the request url, return 400 (Bad Request)
	router.NoRoute(func(c *gin.Context) {
		BadRequest(c, "The endpoint "+c.Request.URL.Path+" is not well formed")
	})

	// Listen and serve on config port
	router.Run(":" + config.Port)
}
