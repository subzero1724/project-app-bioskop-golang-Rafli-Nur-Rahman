package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

// formatValidationError formats validation errors into a readable message
func (v *Validator) formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			messages = append(messages, v.formatFieldError(e))
		}
		return fmt.Errorf(strings.Join(messages, "; "))
	}
	return err
}

// formatFieldError formats a single field error
func (v *Validator) formatFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
