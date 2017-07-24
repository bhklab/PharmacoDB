package api

// Intersection is a combination query model.
type Intersection struct {
	ID      int      `json:"id"`
	ResComb []string `json:"resources_intersected"`
	Path    string   `json:"template-path"`
	Ex      string   `json:"example_path"`
}

// Intersections is a collection of Intersection.
type Intersections []Intersection

// List lists all possible intersections.
func (intersections *Intersections) List() {
	var intersection Intersection
	// First intersection
	intersection.ID = 1
	intersection.ResComb = []string{"cell_line", "drug"}
	intersection.Path = "/intersections/{id}/{cell_id}/{drug_id}"
	intersection.Ex = "https://api.pharmacodb.com/intersections/1/mcf7/paclitaxel?type=name"
	*intersections = append(*intersections, intersection)
	// Second intersection
	intersection.ID = 2
	intersection.ResComb = []string{"cell_line", "dataset"}
	intersection.Path = "/intersections/{id}/{cell_id}/{dataset_id}"
	intersection.Ex = "https://api.pharmacodb.com/intersections/2/mcf7/ccle?type=name"
	*intersections = append(*intersections, intersection)
}

// CellDrugCombination updates receiver with a list of all experiments where a cell line and a drug have been tested.
func (experiments *Experiments) CellDrugCombination(cellID string, drugID string, typ string) error {
	var (
		cell Cell
		drug Drug
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	err = cell.Find(cellID, typ)
	if err != nil {
		return err
	}
	err = drug.Find(drugID, typ)
	if err != nil {
		return err
	}
	query := "SELECT e.experiment_id, t.tissue_id, t.tissue_name, da.dataset_id, da.dataset_name FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.cell_id = ? AND e.drug_id = ?;"
	rows, _ := db.Query(query, cell.ID, drug.ID)
	defer rows.Close()
	for rows.Next() {
		var experiment Experiment
		err = rows.Scan(&experiment.ID, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		experiment.Cell.ID = cell.ID
		experiment.Cell.Name = cell.Name
		experiment.Drug = drug
		err = experiment.DoseResponse()
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*experiments = append(*experiments, experiment)
	}
	return nil
}

// CellDatasetCombination updates receiver with a list of all experiments where a cell line and dataset have been tested.
func (experiments *Experiments) CellDatasetCombination(cellID string, datasetID string, typ string) error {
	var (
		cell    Cell
		dataset Dataset
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	err = cell.Find(cellID, typ)
	if err != nil {
		return err
	}
	err = dataset.Find(datasetID, typ)
	if err != nil {
		return err
	}
	query := "SELECT e.experiment_id, t.tissue_id, t.tissue_name, d.drug_id, d.drug_name FROM experiments e JOIN tissues t ON t.tissue_id = e.tissue_id JOIN drugs d ON d.drug_id = e.drug_id WHERE e.cell_id = ? AND e.dataset_id = ?;"
	rows, _ := db.Query(query, cell.ID, dataset.ID)
	defer rows.Close()
	for rows.Next() {
		var experiment Experiment
		err = rows.Scan(&experiment.ID, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name)
		if err != nil {
			LogPrivateError(err)
			return err
		}
		experiment.Cell.ID = cell.ID
		experiment.Cell.Name = cell.Name
		experiment.Dataset = dataset
		err = experiment.DoseResponse()
		if err != nil {
			LogPrivateError(err)
			return err
		}
		*experiments = append(*experiments, experiment)
	}
	return nil
}
