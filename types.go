package main

import null "gopkg.in/guregu/null.v3"

// Cell type struct
type Cell struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Accession null.String `json:"accession"`
	Tissue    string      `json:"tissue"`
}

// Cells type struct
type Cells struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Accession null.String `json:"accession"`
}
