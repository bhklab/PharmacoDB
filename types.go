package main

type Cell struct {
	ID     int       `json:"id"`
	ACC    *string   `json:"accession_id,omitempty"`
	Name   string    `json:"name"`
	Tissue *Tissue   `json:"tissue,omitempty"`
	SYNS   []Synonym `json:"synonyms,omitempty"`
}

// Tissue is a tissue datatype.
type Tissue struct {
	ID   int       `json:"id"`
	Name *string   `json:"name,omitempty"`
	SYNS []Synonym `json:"synonyms,omitempty"`
}

// Drug is a drug datatype.
type Drug struct {
	ID   int       `json:"id"`
	Name string    `json:"name"`
	SYNS []Synonym `json:"synonyms,omitempty"`
}

// Dataset is a dataset datatype.
type Dataset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Experiment is an experiment datatype.
type Experiment struct {
	ID      int            `json:"experiment_id"`
	Cell    Cell           `json:"cell_line"`
	Tissue  Tissue         `json:"tissue"`
	Drug    Drug           `json:"drug"`
	Dataset Dataset        `json:"dataset"`
	DR      []DoseResponse `json:"dose_responses,omitempty"`
}

// DoseResponse is a dose_response datatype.
type DoseResponse struct {
	Dose     float64 `json:"dose"`
	Response float64 `json:"response"`
}

// Synonym is a match between a datatype name and datasets that use the name.
type Synonym struct {
	Name     string   `json:"name"`
	Datasets []string `json:"datasets"`
}
