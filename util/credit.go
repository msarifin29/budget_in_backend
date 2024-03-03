package util

import "github.com/go-playground/validator/v10"

func IsSupportedTypeCredit(status string) bool {
	switch status {
	case MONTHLY:
		return true
	}
	return false
}

var ValidTypeCredit validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if category, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedTypeCredit(category)
	}
	return false
}

const (
	ACTIVE    = "active"
	COMPLETED = "completed"
)

func IsSupportedStatusHistoryCredit(status string) bool {
	switch status {
	case ACTIVE, COMPLETED:
		return true
	}
	return false
}

var ValidStatusHistoryCredit validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if category, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedStatusHistoryCredit(category)
	}
	return false
}

//  Constants for all supported category credit
const (
	ELECTRONIC     = "electronic"
	SMARTPHONE     = "smathphone"
	LAPTOP         = "laptop"
	COMPUTER       = "computer"
	MOTORCYCLE     = "motorcycle"
	CAR            = "car"
	PROPERTY       = "property"
	FURNITURE      = "furniture"
	KITCHENSET     = "kitchen set"
	VENTURECAPITAL = "venture capital"
)

func IsSupportedCategoryCredit(status string) bool {
	switch status {
	case ELECTRONIC, SMARTPHONE, LAPTOP, COMPUTER, MOTORCYCLE, CAR, PROPERTY, FURNITURE, KITCHENSET, VENTURECAPITAL:
		return true
	}
	return false
}

var ValidCategoryCredit validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if category, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedCategoryCredit(category)
	}
	return false
}

func IsSupportedLoanTermWeekly(n float64) bool {
	switch n {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12:
		return true
	}
	return false
}

var ValidLoanTermWeekly validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if n, ok := fieldLevel.Field().Interface().(float64); ok {
		return IsSupportedLoanTermWeekly(n)
	}
	return false
}

func IsSupportedLoanTermMonthly(n int) bool {
	switch n {
	case 3, 6, 9, 12, 18, 24, 30, 36, 48, 50, 62, 120, 240:
		return true
	}
	return false
}

var ValidLoanTermMonthly validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if n, ok := fieldLevel.Field().Interface().(int); ok {
		return IsSupportedLoanTermMonthly(n)
	}
	return false
}
