package util

import "time"

func IsValidYearMonth(year string, month string) bool {
	// Attempt to parse the input string as a date in YYYY-MM format
	_, err := time.Parse("2006-01", year+"-"+month)
	// If there is no error, the input is valid
	if err == nil {
		return true
	}
	// If there is an error, the input is invalid
	return false
}
