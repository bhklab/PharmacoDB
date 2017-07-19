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
	// GET requests
	Route{GET, "/cell_lines", CellsHandler},
	Route{GET, "/cell_lines/:id", CellHandler},
	Route{GET, "/tissues", TissuesHandler},
	Route{GET, "/tissues/:id", TissueHandler},
	Route{GET, "/drugs", DrugsHandler},
	Route{GET, "/drugs/:id", DrugHandler},
	Route{GET, "/datasets", DatasetsHandler},
	Route{GET, "/datasets/:id", DatasetHandler},
	Route{GET, "/experiments", ExperimentsHandler},
	Route{GET, "/stats", ExperimentsHandler},
}
