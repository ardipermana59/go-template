package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   toSnakeCase(e.Field()),
				Message: formatErrorMessage(e),
			})
		}
	}

	return errors
}

func formatErrorMessage(e validator.FieldError) string {
	field := toSnakeCase(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required", field)
	case "email":
		return fmt.Sprintf("The %s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("The %s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("The %s must not exceed %s characters", field, e.Param())
	case "eqfield":
		return fmt.Sprintf("The %s must match %s", field, toSnakeCase(e.Param()))
	default:
		return fmt.Sprintf("The %s is invalid", field)
	}
}

func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
