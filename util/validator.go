package util

import "github.com/go-playground/validator/v10"

//  Constants for all supported currencies
const (
	IDR = "IDR"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case IDR:
		return true
	}
	return false
}

var ValidCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedCurrency(currency)
	}
	return false
}

//  Constants for all supported user type
const (
	PERSONAL = "personal"
	GROUP    = "group"
)

func IsSupportedType(typeX string) bool {
	switch typeX {
	case PERSONAL, GROUP:
		return true
	}
	return false
}

var ValidType validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if typeX, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedType(typeX)
	}
	return false
}

//  Constants for all supported expense type
const (
	CASH   = "Cash"
	CREDIT = "Credit"
	DEBIT  = "Debit"
)

func IsSupportedExpenseType(typeX string) bool {
	switch typeX {
	case CASH, CREDIT, DEBIT:
		return true
	}
	return false
}

var ValidExpenseType validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if typeX, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedExpenseType(typeX)
	}
	return false
}
