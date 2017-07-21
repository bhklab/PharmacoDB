package models

// Annotation is model for names used by various datasets for each model.
type Annotation struct {
	Name     string   `json:"name"`
	Datasets []string `json:"datasets"`
}

// Annotations is a collection of Annotation.
type Annotations []Annotation
