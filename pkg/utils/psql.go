package utils

import "strings"

// IsDuplicateEntryError checks if the given error is a duplicate entry error.
func IsDuplicateEntryError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

// IsNoRowsError checks if the given error is a no rows error.
func IsNoRowsError(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}
