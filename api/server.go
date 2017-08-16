package api

import "github.com/gin-gonic/gin"

// Init server.
func Init(config *Config) {
	// Set database credentials
	SetDB(config.Version)

	// Gin router with default middleware: logger and recovery
	router := gin.Default()

	v := router.Group(config.Version + "/")
	for _, route := range routes {
		v.Handle(route.Method, route.Endpoint, route.Handler)
	}

	// If no routers match the request url, return 400 (Bad Request)
	router.NoRoute(func(c *gin.Context) {})

	// Listen and serve on config port
	router.Run(":" + config.Port)
}
