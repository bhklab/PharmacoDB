package models

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
