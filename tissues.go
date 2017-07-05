package main

// Tissue is a tissue datatype.
type Tissue struct {
	ID   int     `json:"id"`
	Name *string `json:"name,omitempty"`
}
