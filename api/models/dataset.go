package models

// Dataset is model for a dataset.
type Dataset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Datasets is a collection of Dataset.
type Datasets []Dataset
