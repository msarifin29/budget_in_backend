package util

import "github.com/go-playground/validator/v10"

const (
	OTHER   = "other"
	DAILY   = "daily"
	WEEKLY  = "weekly"
	MONTHLY = "monthly"
)

// Category income
const ( // other id 1
	BUSINES          = "business"          // id 2
	SALARY           = "salary"            // id 3
	ADDITIONALiNCOME = "additional income" // id 4
	LOAN             = "loan"              // id 5
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
	FOODANDDRINK    = "food and drink"    // id 2
	SHOPPING        = "shopping"          // id 3
	TRANSPORT       = "transport"         // id 4
	MOTORCYCLEORCAR = "motorcycle or car" // id 5
	TRAVELING       = "traveling"         // id 6
	HEALTY          = "healty"            // id 7
	COSTANDBILL     = "cost and bill"     // id 8
	EDUCATION       = "education"         // id 9
	SPORTANDHOBBY   = "sport and hobby"   // id 10
	BEAUTY          = "beauty"            // id 11
	WORK            = "work"              // id 12
	FOODINGREDIENTS = "food ingredients"  // id 13
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
