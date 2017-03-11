package main

type Cell struct {
	Id        int            `json:"id"`
	Accession sql.NullString `json:"accession"`
	Name      string         `json:"name"`
	Tissue    sql.NullString `json:"tissue"`
}
