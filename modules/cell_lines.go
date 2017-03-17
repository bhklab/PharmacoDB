package main

import (
	"log"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/guregu/null.v3"
)

// UnfilCell is a list of the same cell line with different synonyms
type UnfilCell struct {
	ID        int
	Name      string
	Accession null.String
	Tissue    Tissue
	Source    string
	SCName    string
}

// GetCLines handles GET requests for cell lines
func GetCLines(c *gin.Context) {
	var (
		cell  Cell
		cells []Cell
	)

	db := InitDb()
	defer db.Close()

	rows, err := db.Query("select cell_id, accession_id, cell_name from cells;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&cell.ID, &cell.Accession, &cell.Name)
		if err != nil {
			log.Fatal(err)
		}
		cells = append(cells, cell)
	}
	defer rows.Close()

	result := gin.H{
		"category":    "cell_lines",
		"description": "list of all cell lines in pharmacodb",
		"count":       len(cells),
		"data":        cells,
	}
	c.IndentedJSON(http.StatusOK, result)
}

// GetCLineByID handles GET request for a cell line using ID
func GetCLineByID(c *gin.Context) {
	var (
		cell     Cell
		ucell    UnfilCell
		synonym  Synonym
		synonyms []Synonym
	)

	db := InitDb()
	defer db.Close()

	cellExists := make(map[string]bool)

	id := c.Param("id")
	queryStr := "select c.cell_id, c.accession_id, c.cell_name, t.tissue_id, t.tissue_name, s.source_name, scn.cell_name from cells c inner join tissues t on t.tissue_id = c.tissue_id inner join source_cell_names scn on scn.cell_id = c.cell_id inner join sources s on s.source_id = scn.source_id where c.cell_id = ?"
	rows, err := db.Query(queryStr, id)
	if err != nil {
		log.Fatal(err)
	}
	j := 0
	for rows.Next() {
		err = rows.Scan(&ucell.ID, &ucell.Accession, &ucell.Name, &ucell.Tissue.ID, &ucell.Tissue.Name, &ucell.Source, &ucell.SCName)
		if err != nil {
			log.Fatal(err)
		}
		if j == 0 {
			cell.ID = ucell.ID
			cell.Accession = ucell.Accession
			cell.Name = ucell.Name
			cell.Tissue.ID = ucell.Tissue.ID
			cell.Tissue.Name = ucell.Tissue.Name
			synonym.Name = ucell.SCName
			synonym.Sources = append(synonym.Sources, ucell.Source)
			synonyms = append(synonyms, synonym)
			j++
			cellExists[ucell.SCName] = true
		} else if cellExists[ucell.SCName] {
			for i, syn := range synonyms {
				if syn.Name == ucell.SCName {
					synonyms[i].Sources = append(synonyms[i].Sources, ucell.Source)
					break
				}
			}
		} else if !cellExists[ucell.SCName] {
			synonym.Name = ucell.SCName
			synonym.Sources = append(synonym.Sources, ucell.Source)
			// synonym.Sources = append(synonym.Sources, ucell.Source)
			synonyms = append(synonyms, synonym)
			cellExists[ucell.SCName] = true
		}
	}
	defer rows.Close()

	cell.Synonyms = synonyms

	if cell.ID == 0 {
		ErrHandler(c, http.StatusNotFound, "Cell line with ID - "+id+" - not found in database.")
	} else {
		c.IndentedJSON(http.StatusOK, cell)
	}
}

// GetCLineByName handles GET request for a cell line using name
func GetCLineByName(c *gin.Context) {
	var cell Cell

	db := InitDb()
	defer db.Close()

	name := c.Param("name")
	row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_name = ?;", name)
	err := row.Scan(&cell.ID, &cell.Accession, &cell.Name, &cell.Tissue)
	if err != nil {
		ErrHandler(c, http.StatusNotFound, "Cell line with name - "+name+" - not found in database.")
	} else {
		c.IndentedJSON(http.StatusOK, cell)
	}
}
