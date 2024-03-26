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
func Date(date string) *time.Time {
	now := time.Now()
	if date == "" {
		return &now
	}
	newDate, _ := time.Parse("2006-01-02", date)
	return &newDate
}

func GetTotalMonth(start *time.Time, end *time.Time) (totalMonths int) {
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

	// Calculate the difference in years
	yearsDiff := endDate.Year() - startDate.Year()

	// Calculate the difference in months (considering edge cases)
	monthsDiff := int(endDate.Month() - startDate.Month())
	if endDate.Day() < startDate.Day() {
		monthsDiff-- // Adjust for cases where end date falls on a day before start date within the same month
	}

	// Handle negative month difference (start date after end date)
	if monthsDiff < 0 {
		monthsDiff += 12 // Add 12 to account for wrapping around to the previous year
		yearsDiff--      // Decrement years to compensate for the additional year
	}

	// Calculate total months
	totalMonths = yearsDiff*12 + monthsDiff
	return totalMonths
}
func GenerateDates(start time.Time, end time.Time) []string {

	// Ensure start date is before or equal to end date
	if start.After(end) {
		return []string{} // Return empty slice if start date is after end date
	}

	var dates []string
	for currentDate := start; !currentDate.After(end); currentDate = currentDate.AddDate(0, 1, 0) {
		// Adjust for end date being February 29th in a leap year
		if currentDate.Equal(end) && end.Month() == time.February && end.Day() == 29 && !isLeapYear(currentDate.Year()) {
			currentDate = currentDate.AddDate(0, 0, -1) // Adjust to last day of January if not a leap year
		}
		// Preserve the day of the month from the start date
		currentDate = time.Date(currentDate.Year(), currentDate.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

		dates = append(dates, currentDate.Format("2006-01-02"))
	}
	return dates
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}
