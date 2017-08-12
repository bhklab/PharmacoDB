package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

	// Rate limit
	// limiter := tollbooth.NewLimiter(1, time.Second)
	//
	// // Limit only GET and POST requests.
	// limiter.Methods = []string{"GET", "POST"}

	router := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.StaticFile("/favicon.ico", "./static/images/favicon.png")

	v := router.Group(Version() + "/")
	for _, route := range routes {
		v.Handle(route.Method, route.Endpoint, route.Handler)
	}

	// Respond with status code 400 (Bad Request) if no routers match the request url.
	router.NoRoute(func(c *gin.Context) { LogBadRequestError(c) })

	router.Run(":" + Port())
}
