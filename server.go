package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Init API server.
// Using the gin router.
func Init() {
	router := gin.Default()

	// Use appropriate mode for router.
	if Mode() == DebugMode {
		gin.SetMode(gin.DebugMode)
	}
	if Mode() == ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	if Mode() == TestMode {
		gin.SetMode(gin.TestMode)
	}

	v := router.Group(APIVersion())

	for _, route := range routes {
		v.Handle(route.Method, route.Endpoint, route.Handler)
	}

	// Respond with status code 400 (Bad Request)
	// if no routers match the request url.
	router.NoRoute(func(c *gin.Context) {
		LogPublicError(c, ErrorTypePublic, http.StatusBadRequest, "Bad Request")
	})

	router.Run(":" + Port())
}
