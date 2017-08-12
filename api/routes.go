package api

import "github.com/gin-gonic/gin"

// HTTP methods.
const (
	GET    string = "GET"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"
	HEAD   string = "HEAD"
	OPTION string = "OPTION"
	PATCH  string = "PATCH"
)

// Route is a routing model.
type Route struct {
	Method   string
	Endpoint string
	Handler  gin.HandlerFunc
}

// Routes is a collection of Route.
type Routes []Route

var routes = Routes{
	Route{GET, "/cell_lines", ListCellsHandler},
	Route{GET, "/tissues", ListCellsHandler},
	Route{GET, "/compounds", ListCellsHandler},
	Route{GET, "/datasets", ListCellsHandler},
	Route{GET, "/experiments", ListCellsHandler},
}
