package util

import "github.com/go-playground/validator/v10"

const (
	OTHER   = "Other"
	DAILY   = "daily"
	WEEKLY  = "weekly"
	MONTHLY = "monthly"
)

// Category income
const ( // other id 1
	BUSINES          = "Business"          // id 2
	SALARY           = "Salary"            // id 3
	ADDITIONALiNCOME = "Additional Income" // id 4
	LOAN             = "Loan"              // id 5
)

func IsSupportedCategoryType(status string) bool {
	switch status {
	case OTHER, BUSINES, SALARY, ADDITIONALiNCOME, LOAN:
		return true
	}
	return false
}

var ValidCategoryType validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if category, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedCategoryType(category)
	}
	return false
}

func InputCategoryIncome(id float64) string {
	switch id {
	case 1:
		return OTHER
	case 2:
		return BUSINES
	case 3:
		return SALARY
	case 4:
		return ADDITIONALiNCOME
	case 5:
		return LOAN
	default:
		return OTHER
	}
}

// func IsSupportedCategoryIncome(status string) bool {
// 	switch status {
// 	case OTHER, DAILY, WEEKLY, MONTHLY:
// 		return true
// 	}
// 	return false
// }

// var ValidCategoryIncome validator.Func = func(fieldLevel validator.FieldLevel) bool {
// 	if category, ok := fieldLevel.Field().Interface().(string); ok {
// 		return IsSupportedCategoryIncome(category)
// 	}
// 	return false
// }

// Category expense
const ( // other id 1
	FOODANDDRINK    = "Food and Drink"    // id 2
	SHOPPING        = "Shopping"          // id 3
	TRANSPORT       = "Transport"         // id 4
	MOTORCYCLEORCAR = "Motorcycle or Car" // id 5
	TRAVELING       = "Traveling"         // id 6
	HEALTY          = "Healty"            // id 7
	COSTANDBILL     = "Cost and Bill"     // id 8
	EDUCATION       = "Education"         // id 9
	SPORTANDHOBBY   = "Sport and Hobby"   // id 10
	BEAUTY          = "Beauty"            // id 11
	WORK            = "Work"              // id 12
	FOODINGREDIENTS = "Food Ingredients"  // id 13
)

func IsSupportedCategoryExpense(status string) bool {
	switch status {
	case OTHER, FOODANDDRINK, SHOPPING, TRANSPORT, MOTORCYCLEORCAR, TRAVELING, HEALTY, COSTANDBILL, EDUCATION, SPORTANDHOBBY, BEAUTY, WORK, FOODINGREDIENTS:
		return true
	}
	return false
}

var ValidCategoryExpense validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if category, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedCategoryExpense(category)
	}
	return false
}

func InputCategoryexpense(id float64) string {
	switch id {
	case 1:
		return OTHER
	case 2:
		return FOODANDDRINK
	case 3:
		return SHOPPING
	case 4:
		return TRANSPORT
	case 5:
		return MOTORCYCLEORCAR
	case 6:
		return TRAVELING
	case 7:
		return HEALTY
	case 8:
		return COSTANDBILL
	case 9:
		return EDUCATION
	case 10:
		return SPORTANDHOBBY
	case 11:
		return BEAUTY
	case 12:
		return WORK
	case 13:
		return FOODINGREDIENTS
	default:
		return OTHER
	}
}
