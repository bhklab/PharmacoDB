package main

// Cell is model for a cell line.
type Cell struct {
	ID          int         `json:"id"`
	ACC         *string     `json:"accession_id,omitempty"`
	Name        string      `json:"name"`
	Tissue      *Tissue     `json:"tissue,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

// Cells is a collection of Cell.
type Cells []Cell

// Tissue is model for a tissue.
type Tissue struct {
	ID          int         `json:"id"`
	Name        *string     `json:"name,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

// Tissues is a collection of Tissue.
type Tissues []Tissue

// Drug is model for a drug.
type Drug struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Annotations Annotations `json:"annotations,omitempty"`
}

// Drugs is a collection of Drug.
type Drugs []Drug

// Dataset is model for a dataset.
type Dataset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Datasets is a collection of Dataset.
type Datasets []Dataset

// Experiment is model for an experiment.
type Experiment struct {
	ID            int           `json:"experiment_id"`
	Cell          Cell          `json:"cell_line"`
	Tissue        Tissue        `json:"tissue"`
	Drug          Drug          `json:"drug"`
	Dataset       Dataset       `json:"dataset"`
	DoseResponses DoseResponses `json:"dose_responses,omitempty"`
}

// Experiments is a collection of Experiment.
type Experiments []Experiment

// DoseResponse is model for dose/response data.
type DoseResponse struct {
	Dose     float64 `json:"dose"`
	Response float64 `json:"response"`
}

// DoseResponses is a collection of DoseResponse.
type DoseResponses []DoseResponse

// Annotation is model for names used by various datasets for each model.
type Annotation struct {
	Name     string   `json:"name"`
	Datasets []string `json:"datasets"`
}

// Annotations is a collection of Annotation.
type Annotations []Annotation
