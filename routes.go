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
	// Tissues routes
	Route{GET, "/tissues", TissuesHandler},
	// Drugs routes
	Route{GET, "/drugs", DrugsHandler},
	// Datasets routes
	Route{GET, "/datasets", DatasetsHandler},
	// Experiments routes
	Route{GET, "/experiments", ExperimentsHandler},
}
