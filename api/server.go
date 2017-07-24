package api

import "github.com/gin-gonic/gin"

// Context configuration.
type Context struct {
	Mode    string
	Port    string
	Version string
}

// Init new server, using gin router.
func Init(c *Context) {
	SetMode(c.Mode)
	SetPort(c.Port)
	SetVersion(c.Version)

	SetDB()

	router := gin.Default()

	router.StaticFile("/favicon.ico", "./static/images/favicon.png")

	v := router.Group(Version())
	for _, route := range routes {
		v.Handle(route.Method, route.Endpoint, route.Handler)
	}

	// Respond with status code 400 (Bad Request) if no routers match the request url.
	router.NoRoute(func(c *gin.Context) { LogBadRequestError(c) })

	router.Run(":" + Port())
}
