package api

import "fmt"

// Count returns the total number of records in table.
func Count(table string) (int, error) {
	var count int
	db, err := Database()
	defer db.Close()
	if err != nil {
		return count, err
	}
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s;", table)).Scan(&count)
	if err != nil {
		LogSentry(err)
	}
	return count, err
}

// isEmptyStr returns true if string is empty, and false otherwise.
func isEmptyStr(s string) bool {
	return s == ""
}

// isSameStr compares two strings, and returns true if they are the same. Otherwise, returns false.
func isSameStr(a string, b string) bool {
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
