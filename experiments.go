package main

import (
	"database/sql"
	"fmt"
)

// PaginatedExperiments returns a list of paginated experiments.
func PaginatedExperiments(page int, limit int) (Experiments, error) {
	var (
		experiment  Experiment
		experiments Experiments
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	s := (page - 1) * limit
	query := fmt.Sprintf("SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id limit %d,%d;", s, limit)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return nil, err
		}
		experiments = append(experiments, experiment)
	}
	return experiments, nil
}

// FindExperiment returns an experiment, selected using its ID.
func FindExperiment(id string) (Experiment, error) {
	var experiment Experiment
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return experiment, err
	}
	query := "SELECT e.experiment_id, c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, d.drug_id, d.drug_name, da.dataset_id, da.dataset_name FROM experiments e JOIN cells c ON c.cell_id = e.cell_id JOIN tissues t ON t.tissue_id = e.tissue_id JOIN drugs d ON d.drug_id = e.drug_id JOIN datasets da ON da.dataset_id = e.dataset_id WHERE e.experiment_id = ?;"
	row := db.QueryRow(query, id)
	err = row.Scan(&experiment.ID, &experiment.Cell.ID, &experiment.Cell.Name, &experiment.Tissue.ID, &experiment.Tissue.Name, &experiment.Drug.ID, &experiment.Drug.Name, &experiment.Dataset.ID, &experiment.Dataset.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			LogPrivateError(ErrorTypePrivate, err)
		}
		return experiment, err
	}
	return experiment, nil
}

// DoseResponse retrieves dose/response data for an experiment.
func (experiment *Experiment) DoseResponse() error {
	var doseResponse DoseResponse
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT dose, response FROM dose_responses WHERE experiment_id = ?;", experiment.ID)
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&doseResponse.Dose, &doseResponse.Response)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return err
		}
		experiment.DoseResponses = append(experiment.DoseResponses, doseResponse)
	}
	return nil
}
