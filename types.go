package main

import "gopkg.in/guregu/null.v3"

type (
	// SynonymReduced is a synonym match between a name and a single dataset.
	SynonymReduced struct {
		Name    string `json:"name"`
		Dataset string `json:"dataset"`
	}

	// Synonym is a synonym match between a name and multiple datasets.
	Synonym struct {
		Name     string   `json:"name"`
		Datasets []string `json:"datasets"`
	}

	// DatasetStat contains the number of a resource tested in a dataset.
	DatasetStat struct {
		Dataset int `json:"dataset"`
		Count   int `json:"count"`
	}
)

type (
	// CellReduced is a cell line with only ID and Name attributes.
	CellReduced struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// Cell is a datatype.
	Cell struct {
		ID        int           `json:"id"`
		Accession null.String   `json:"accession"`
		Name      string        `json:"name"`
		Tissue    TissueReduced `json:"tissue"`
		Synonyms  []Synonym     `json:"synonyms"`
	}
)

type (
	// TissueReduced is a tissue with only ID and Name attributes.
	TissueReduced struct {
		ID   int         `json:"id"`
		Name null.String `json:"name"`
	}

	// Tissue is a datatype.
	Tissue struct {
		ID   int         `json:"id"`
		Name null.String `json:"name"`
	}
)

type (
	// DrugReduced is a drug with only ID and Name attributes.
	DrugReduced struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// Drug is a datatype.
	Drug struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

type (
	// DatasetReduced is a dataset with only ID and Name attributes.
	DatasetReduced struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// Dataset is a datatype.
	Dataset struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)
