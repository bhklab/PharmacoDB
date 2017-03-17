package main

import null "gopkg.in/guregu/null.v3"

// Tissue is a tissue datatype
type Tissue struct {
	ID   int         `json:"id"`
	Name null.String `json:"name"`
}
