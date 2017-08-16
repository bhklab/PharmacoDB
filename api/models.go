package api

// Cell is a resource model for cell_lines.
type Cell struct {
	ID          int         `json:"id"`
	ACC         *string     `json:"accession_id,omitempty"`
	Name        string      `json:"name"`
	Tissue      *Tissue     `json:"tissue,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

// Cells is a collection of Cell.
type Cells []Cell

// Tissue is a resource model for tissues.
type Tissue struct {
	ID          int         `json:"id"`
	Name        *string     `json:"name,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

// Tissues is a collection of Tissue.
type Tissues []Tissue

// Compound is a resource model for compounds.
type Compound struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Annotations Annotations `json:"annotations,omitempty"`
}

// Compounds is a collection of Drug.
type Compounds []Compound

// Dataset is a resource model for datasets.
type Dataset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Datasets is a collection of Dataset.
type Datasets []Dataset

// Experiment is a resource model for experiments.
type Experiment struct {
	ID            int           `json:"experiment_id"`
	Cell          Cell          `json:"cell_line"`
	Tissue        Tissue        `json:"tissue"`
	Compound      Compound      `json:"compound"`
	Dataset       Dataset       `json:"dataset"`
	DoseResponses DoseResponses `json:"dose_responses,omitempty"`
}

// Experiments is a collection of Experiment.
type Experiments []Experiment

// DoseResponse is a model for a dose/response data pair.
type DoseResponse struct {
	Dose     float64 `json:"dose"`
	Response float64 `json:"response"`
}

// DoseResponses is a collection of DoseResponse.
type DoseResponses []DoseResponse

// Annotation models the name used by datasets for a resource item.
type Annotation struct {
	Name     string   `json:"name"`
	Datasets []string `json:"datasets"`
}

// Annotations is a collection of Annotation.
type Annotations []Annotation
