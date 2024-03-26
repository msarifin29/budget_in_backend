package util

import (
	"fmt"
	"time"
)

func GeneratedMonth(month int) (m time.Time) {
	// Set today's date
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	fmt.Scanf("%d", &month)
	for i := 0; i < month; i++ {
		nextMonth := today.AddDate(0, i+1, 0) // Add i+1 months because i starts from 0

		// Adjust for edge case: February might have less than 19 days
		if nextMonth.Month() == time.February && nextMonth.Day() > 28 {
			nextMonth = nextMonth.AddDate(0, 0, -1*(nextMonth.Day()-28)) // Adjust to the last day of January if needed
		}
		m = nextMonth
	}
	return m
}
func GetMonth(month *time.Time) int {
	currentMonth := month.Month()
	var monthInt int
	switch currentMonth {
	case time.January:
		monthInt = 1
	case time.February:
		monthInt = 2
	case time.March:
		monthInt = 3
	case time.April:
		monthInt = 4
	case time.May:
		monthInt = 5
	case time.June:
		monthInt = 6
	case time.July:
		monthInt = 7
	case time.August:
		monthInt = 8
	case time.September:
		monthInt = 9
	case time.October:
		monthInt = 10
	case time.November:
		monthInt = 11
	case time.December:
		monthInt = 12
	}
	return monthInt
}
