package models

// Cell is model for a cell line.
type Cell struct {
	ID          int         `json:"id"`
	ACC         string      `json:"accession_id,omitempty"`
	Name        string      `json:"name"`
	Tissue      Tissue      `json:"tissue,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

// Cells is a collection of Cell.
type Cells []Cell
