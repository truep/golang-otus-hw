package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrValidateLen    = errors.New("value does not match required length")
	ErrValidateRegexp = errors.New("value does not match regexp rule")
	ErrValidateIn     = errors.New("value is not in set")
	ErrValidateMin    = errors.New("value is smaller than min rule")
	ErrValidateMax    = errors.New("value is larger than max rule")

	ErrUnsupportedType  = errors.New("unsupported type, expected struct")
	ErrUnsupportedField = errors.New("unsupported field type for validate")
	ErrInvalidRule      = errors.New("invalid rule")
	ErrInvalidRuleValue = errors.New("invalid rule value")
)

type ValidationError struct {
	Field string
	Err   error
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Err)
}

func (ve ValidationError) Unwrap() error {
	return ve.Err
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("Validation errors:\n")

	for _, item := range v {
		b.WriteString(item.Error() + "\n")
	}
	return b.String()
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	var t reflect.Type

	switch {
	case value.Kind() == reflect.Ptr:
		t = value.Type().Elem()
		if t.Kind() == reflect.Struct {
			value = reflect.Indirect(value)
		}
	case value.Kind() == reflect.Struct:
		t = reflect.TypeOf(v)
	default:
		return ErrUnsupportedType
	}

	errorList := make(ValidationErrors, 0, value.NumField())
	var ve *ValidationErrors
	err := validateFields(value, t, &errorList)
	if err != nil && !errors.As(err, &ve) {
		return err
	}

	if len(errorList) > 0 {
		return errorList
	}

	return nil
}

func validateFields(value reflect.Value, t reflect.Type, errorList *ValidationErrors) error {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := t.Field(i)

		var ve *ValidationErrors
		switch {
		case field.Kind() == reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				fieldValue := field.Index(j)
				err := validateField(fieldType, fieldValue, errorList)
				if err != nil && !errors.As(err, &ve) {
					return err
				}
			}
		default:
			err := validateField(fieldType, field, errorList)
			if err != nil && !errors.As(err, &ve) {
				return err
			}
		}
	}

	return errorList
}

func validateField(fieldType reflect.StructField, fieldValue reflect.Value, errorList *ValidationErrors) error {
	validateRules := getFieldValidationRules(fieldType)

	for _, rule := range validateRules {
		if err := validate(rule, fieldType.Name, fieldValue); err != nil {
			if errors.As(err, &ValidationError{}) {
				*errorList = append(*errorList, ValidationError{Field: fieldType.Name, Err: errors.Unwrap(err)})
			} else {
				return err
			}
		}
	}

	return errorList
}

func getFieldValidationRules(field reflect.StructField) []string {
	validateTag, ok := field.Tag.Lookup("validate")
	if !ok {
		return []string{}
	}
	return strings.Split(validateTag, "|")
}

func validate(rule, fieldName string, value reflect.Value) error {
	rules := strings.Split(rule, ":")
	if len(rules) != 2 {
		return ErrInvalidRule
	}

	ruleName := rules[0]
	ruleValue := rules[1]

	if ruleValue == "" {
		return ErrInvalidRuleValue
	}

	switch {
	case value.Kind() == reflect.String:
		switch ruleName {
		case "len":
			return validateLen(ruleValue, fieldName, value)
		case "regexp":
			return validateRegexp(ruleValue, fieldName, value)
		case "in":
			return validateIn(ruleValue, fieldName, value)
		default:
			return ErrInvalidRule
		}
	case value.Kind() == reflect.Int:
		switch ruleName {
		case "in":
			return validateIn(ruleValue, fieldName, value)
		case "min":
			return validateMin(ruleValue, fieldName, value)
		case "max":
			return validateMax(ruleValue, fieldName, value)
		default:
			return ErrInvalidRule
		}
	default:
		return ErrUnsupportedField
	}
}

func validateLen(ruleValue, fieldName string, value reflect.Value) error {
	length, err := strconv.Atoi(ruleValue)
	if err != nil {
		return ErrInvalidRuleValue
	}
	if value.Len() != length {
		return ValidationError{fieldName, ErrValidateLen}
	}

	return nil
}

func validateRegexp(ruleValue, fieldName string, value reflect.Value) error {
	matched, err := regexp.MatchString(ruleValue, value.String())
	if err != nil {
		return ErrInvalidRuleValue
	}
	if !matched {
		return ValidationError{fieldName, ErrValidateRegexp}
	}

	return nil
}

func validateIn(ruleValue, fieldName string, value reflect.Value) error {
	validValues := strings.Split(ruleValue, ",")
	if value.Kind() == reflect.String {
		if !stringInSlice(value.String(), validValues) {
			return ValidationError{fieldName, ErrValidateIn}
		}
	} else if value.Kind() == reflect.Int {
		intVal := strconv.Itoa(int(value.Int()))
		if !stringInSlice(intVal, validValues) {
			return ValidationError{fieldName, ErrValidateIn}
		}
	}

	return nil
}

func validateMin(ruleValue, fieldName string, value reflect.Value) error {
	minVal, err := strconv.Atoi(ruleValue)
	if err != nil {
		return ErrInvalidRuleValue
	}
	if int(value.Int()) < minVal {
		return ValidationError{fieldName, ErrValidateMin}
	}

	return nil
}

func validateMax(ruleValue, fieldName string, value reflect.Value) error {
	maxVal, err := strconv.Atoi(ruleValue)
	if err != nil {
		return ErrInvalidRuleValue
	}
	if int(value.Int()) > maxVal {
		return ValidationError{fieldName, ErrValidateMax}
	}

	return nil
}

func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
