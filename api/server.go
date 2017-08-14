package api

import "github.com/gin-gonic/gin"

// Init new server.
func Init(config *Config) {
	// Gin router with default middleware: logger and recovery
	router := gin.Default()

	// Set database credentials
	SetDB(config.Version)

	v := router.Group(config.Version + "/")
	for _, route := range routes {
		v.Handle(route.Method, route.Endpoint, route.Handler)
	}

	// If no routers match the request url, return code 400 (Bad Request)
	router.NoRoute(func(c *gin.Context) {})

	// Listen and serve on config port
	router.Run(":" + config.Port)
}
