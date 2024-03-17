package util

import "time"

func CreatedAt(date string) *time.Time {
	now := time.Now()
	if date == "" {
		return &now
	}
	newDate, _ := time.Parse(time.RFC3339, date)
	return &newDate
}
