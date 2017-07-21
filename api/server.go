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

	v := router.Group(Version())
	for _, route := range routes {
		v.Handle(route.Method, route.Endpoint, route.Handler)
	}

	router.Run(":" + Port())
}
