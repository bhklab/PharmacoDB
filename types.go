package main

import "gopkg.in/guregu/null.v3"

type (
	// Synonym is a synonym match between a name and multiple datasets.
	Synonym struct {
		Name     string   `json:"name"`
		Datasets []string `json:"datasets"`
	}

	// SynonymReduced is a synonym match between a name and a single dataset.
	SynonymReduced struct {
		Name    string `json:"name"`
		Dataset string `json:"dataset"`
	}

	// DataTypeReduced is a datatype with only ID and Name attributes.
	DataTypeReduced struct {
		ID   int         `json:"id"`
		Name null.String `json:"name"`
	}

	// DatasetStat contains the number of a resource tested in a dataset.
	DatasetStat struct {
		Dataset int `json:"dataset"`
		Count   int `json:"count"`
	}
)

// Cell is a cell line datatype.
type Cell struct {
	ID        int             `json:"id"`
	Accession null.String     `json:"accession"`
	Name      string          `json:"name"`
	Tissue    DataTypeReduced `json:"tissue"`
	Synonyms  []Synonym       `json:"synonyms"`
}

// Tissue is a tissue datatype.
type Tissue struct {
	ID   int         `json:"id"`
	Name null.String `json:"name"`
}

// Drug is a drug datatype.
type Drug struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Dataset is a dataset datatype.
type Dataset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
