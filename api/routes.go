package api

import "github.com/gin-gonic/gin"

// HTTP methods.
const (
	GET    string = "GET"
	HEAD   string = "HEAD"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"
	OPTION string = "OPTION"
	PATCH  string = "PATCH"
)

// Route is a routing model.
type Route struct {
	Endpoint string
	Handler  gin.HandlerFunc
}

// Routes is a collection of Route.
type Routes []Route

var routesGET = Routes{
	Route{"/cell_lines", IndexCells},
	Route{"/cell_lines/:id", ShowCell},
	Route{"/cell_lines/:id/compounds", CellCompounds},

	Route{"/tissues", IndexTissues},
	Route{"/tissues/:id", ShowTissue},
	Route{"/tissues/:id/cell_lines", TissueCells},
	Route{"/tissues/:id/compounds", TissueCompounds},

	Route{"/compounds", IndexCompounds},
	Route{"/compounds/:id", ShowCompound},
	Route{"/compounds/:id/cell_lines", CompoundCells},
	Route{"/compounds/:id/tissues", CompoundTissues},

	Route{"/datasets", IndexDatasets},
	Route{"/datasets/:id", ShowDataset},
	Route{"/datasets/:id/cell_lines", DatasetCells},
	Route{"/datasets/:id/tissues", DatasetTissues},
	Route{"/datasets/:id/compounds", DatasetCompounds},

	Route{"/experiments", IndexExperiments},
	Route{"/experiments/:id", ShowExperiment},

	Route{"/intersections", IndexIntersections},
	Route{"/intersections/1/:cell_id/:compound_id", CellCompoundIntersection},
	Route{"/intersections/2/:cell_id/:dataset_id", CellDatasetIntersection},

	Route{"/stats/tissues/cell_lines", StatTissuesCells},

	Route{"/stats/datasets/cell_lines", StatDatasetsCells},
	Route{"/stats/datasets/tissues", StatDatasetsTissues},
	Route{"/stats/datasets/compounds", StatDatasetsCompounds},

	Route{"/stats/datasets/cell_lines/tissues/:id", StatDatasetsTissuesCells},
	Route{"/stats/datasets/cell_lines/compounds/:id", StatDatasetsCompoundsCells},
	Route{"/stats/datasets/tissues/compounds/:id", StatDatasetsCompoundsTissues},
	Route{"/stats/datasets/compounds/cell_lines/:id", StatDatasetsCellsCompounds},
	Route{"/stats/datasets/compounds/tissues/:id", StatDatasetsTissuesCompounds},

	Route{"/stats/datasets/experiments", StatDatasetsExperiments},
}

var routesHEAD = Routes{
	Route{"/cell_lines", IndexCellsHEAD},
}
