package models

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
