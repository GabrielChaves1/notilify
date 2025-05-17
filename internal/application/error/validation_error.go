package errors

import "fmt"

type ValidationRule string

const (
	Required      ValidationRule = "required"
	MinLength     ValidationRule = "min_length"
	MaxLength     ValidationRule = "max_length"
	InvalidFormat ValidationRule = "invalid_format"
	InvalidValue  ValidationRule = "invalid_value"
)

type Validation struct {
	Field      string
	Value      interface{}
	Rule       ValidationRule
	Constraint interface{}
}

func NewValidation(field string, value interface{}, rule ValidationRule, constraint interface{}) Validation {
	return Validation{
		Field:      field,
		Value:      value,
		Rule:       rule,
		Constraint: constraint,
	}
}

func (v Validation) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s (value: %v, constraint: %v)", v.Field, v.Rule, v.Value, v.Constraint)
}

type ValidationErrors []Validation

func (ve *ValidationErrors) Error() string {
	if len(*ve) == 0 {
		return ""
	}

	var errorMessages []string
	for _, v := range *ve {
		errorMessages = append(errorMessages, v.Error())
	}
	return fmt.Sprintf("multiple validation errors: %s", errorMessages)
}

func (ve *ValidationErrors) HasErrors() bool {
	return len(*ve) > 0
}

func (ve *ValidationErrors) Add(field string, value interface{}, rule ValidationRule, constraint interface{}) {
	*ve = append(*ve, NewValidation(field, value, rule, constraint))
}

func (ve *ValidationErrors) Merge(errs *ValidationErrors) {
	if errs == nil {
		return
	}

	*ve = append(*ve, *errs...)
}
