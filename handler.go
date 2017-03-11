package main

import (
	"log"
)

var (
	cell   Cell
	result gin.H
)

func getCell(c *gin.Context) {
	id := c.Param("id")
	row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
	err = row.Scan(&cell.Id, &cell.Accession, &cell.Name, &cell.Tissue)
	if err != nil {
		row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
		err = row.Scan(&cell.Id, &cell.Accession, &cell.Name, &cell.Tissue)
		if err != nil {
			row := db.QueryRow("select cell_id, accession_id, cell_name, tissue_name from cells inner join tissues on cells.tissue_id = tissues.tissue_id where cells.cell_id = ?;", id)
			err = row.Scan(&cell.Id, &cell.Accession, &cell.Name, &cell.Tissue)
			if err != nil {
				result = gin.H{
					"category": "cell line",
					"count":    0,
					"data":     nil,
				}
			} else {
				result = gin.H{
					"category": "cell line",
					"count":    1,
					"data":     cell,
				}
			}
		} else {
			result = gin.H{
				"category": "cell line",
				"count":    1,
				"data":     cell,
			}
		}
	} else {
		result = gin.H{
			"category": "cell line",
			"count":    1,
			"data":     cell,
		}
	}
	c.JSON(http.StatusOK, result)
}
