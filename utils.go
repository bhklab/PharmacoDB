package main

import "reflect"

// sameString returns true if a == b, and false otherwise.
func sameString(a string, b string) bool {
	return a == b
}

// stringInSlice returns true if list contains a string, and false otherwise.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// DataTypeInterface is an interface struct for all db datatypes.
type DataTypeInterface struct {
	obj interface{}
}

// Annotates adds annotation to an interface to any data type.
// [ Mock abstraction model ]
func (obj DataTypeInterface) Annotates(query string) error {
	var (
		annotation     Annotation
		annotations    Annotations
		annotationName string
		datasetName    string
	)
	db, err := InitDB()
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query(query, reflect.ValueOf(obj).FieldByName("ID"))
	defer rows.Close()
	if err != nil {
		LogPrivateError(ErrorTypePrivate, err)
		return err
	}
	exists := make(map[string]bool)
	for rows.Next() {
		err = rows.Scan(&annotationName, &datasetName)
		if err != nil {
			LogPrivateError(ErrorTypePrivate, err)
			return err
		}
		if exists[annotationName] {
			for i, a := range annotations {
				if a.Name == annotationName && !stringInSlice(datasetName, a.Datasets) {
					annotations[i].Datasets = append(annotations[i].Datasets, datasetName)
				}
			}
		} else {
			var datasetsNew []string
			annotation.Name = annotationName
			annotation.Datasets = append(datasetsNew, datasetName)
			annotations = append(annotations, annotation)
			exists[annotationName] = true
		}
	}
	reflect.ValueOf(obj).FieldByName("Annotations").Elem().Set(reflect.ValueOf(annotations))
	return nil
}
