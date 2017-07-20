package models

// Tissue is model for a tissue.
type Tissue struct {
	ID          int         `json:"id"`
	Name        *string     `json:"name,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

// Tissues is a collection of Tissue.
type Tissues []Tissue
