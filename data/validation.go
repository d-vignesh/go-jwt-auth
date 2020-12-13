package data

import (
	"fmt"

	"github.com/go-playground/validator"
)

// ValidationError wraps the validator FieldError so we do not
// expose this to outside code
type ValidationError struct {
	validator.FieldError
}

// Error provides the string format of the validation error
func (v ValidationError) Error() string {
	if v.Tag() == "required" {
		return fmt.Sprintf("%s is required", v.Field())
	}

	return fmt.Sprintf(
		"key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag())
}

// ValidationErrors is a wrapper for list of ValidationError
type ValidationErrors []ValidationError

// Errors convert the ValidationErrors slice into string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

// Validation is the type for validator
type Validation struct {
	validate *validator.Validate
}

// NewValidation returns a Validator instance
func NewValidation() *Validation {
	validate := validator.New()
	return &Validation{validate}
}

// Validate method validates the given struct based on the validate tags
// and returns validation error if any
func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i)
	if errs == nil {
		return nil
	}

	var returnErrs ValidationErrors
	for _, err := range errs.(validator.ValidationErrors) {
		// cast the FieldError into our ValidationError and append to the slice
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}
	return returnErrs
}
