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
	Route{GET, "/cell_lines", CellLinesHandler},
	Route{GET, "/cell_lines/:id", CellLinesHandler},
	// Tissues routes
	Route{GET, "/Tissues", CellLinesHandler},
	Route{GET, "/Tissues/:id", CellLinesHandler},
	// Drugs routes
	Route{GET, "/Drugs", CellLinesHandler},
	Route{GET, "/Drugs/:id", CellLinesHandler},
	// Datasets routes
	Route{GET, "/Datasets", CellLinesHandler},
	Route{GET, "/Datasets/:id", CellLinesHandler},
	// Experiments routes
	Route{GET, "/Experiments", CellLinesHandler},
	Route{GET, "/Experiments/:id", CellLinesHandler},
}
