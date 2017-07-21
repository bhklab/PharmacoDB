package models

// Drug is model for a drug.
type Drug struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Annotations Annotations `json:"annotations,omitempty"`
}

// Drugs is a collection of Drug.
type Drugs []Drug
