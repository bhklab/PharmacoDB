package main

import "github.com/gin-gonic/gin"

// HTTP Request Methods.
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
	// Cell lines routes
	Route{GET, "/cell_lines", CellsHandler},
	Route{GET, "/cell_lines/:id", CellsHandler},
	// Tissues routes
	Route{GET, "/tissues", TissuesHandler},
	Route{GET, "/tissues/:id", CellsHandler},
	// Drugs routes
	Route{GET, "/drugs", CellsHandler},
	Route{GET, "/drugs/:id", CellsHandler},
	// Datasets routes
	Route{GET, "/datasets", CellsHandler},
	Route{GET, "/datasets/:id", CellsHandler},
	// Experiments routes
	Route{GET, "/experiments", CellsHandler},
	Route{GET, "/experiments/:id", CellsHandler},
}
