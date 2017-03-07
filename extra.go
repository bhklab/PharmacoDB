// migration for creating cells table (use as template to create further migrations)
func createCells() {
	stmt, err := db.Prepare("CREATE TABLE cells (cell_id int NOT NULL, accession_id VARCHAR(40), cell_name TEXT, PRIMARY KEY (cell_id));")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("'cells' table successfully migrated ...")
	}
}
