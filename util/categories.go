package util

import "github.com/go-playground/validator/v10"

const (
	OTHER = "other"
)

func IsSupportedCategoryType(status string) bool {
	switch status {
	case OTHER:
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
