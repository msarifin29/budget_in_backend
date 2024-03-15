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
const ( // other id 1
	ELECTRONIC        = "electronic"          // id 2
	HANDTPHONE        = "handphone"           // id 3
	COMPUTERANDLAPTOP = "computer and laptop" // id 4
	MOTORCYCLE        = "motorcycle"          // id 5
	CAR               = "car"                 // id 6
	PROPERTY          = "property"            // id 7
	FURNITURE         = "furniture"           // id 8
	KITCHENSET        = "kitchen set"         // id 9
	VENTURECAPITAL    = "venture capital"     // id 10
)

func IsSupportedCategoryCredit(status string) bool {
	switch status {
	case ELECTRONIC, HANDTPHONE, COMPUTERANDLAPTOP, MOTORCYCLE, CAR, PROPERTY, FURNITURE, KITCHENSET, VENTURECAPITAL, OTHER:
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
