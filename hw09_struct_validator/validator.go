package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	UnsupportedInputType = errors.New("unsupported input type")
)

const (
	validationTag string = "validate"
	tagDivider    string = "|"
	tagDefinder   string = ":"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result string
	result += fmt.Sprintln("errors while validation:")
	for _, err := range v {
		result += fmt.Sprintf("%v\t%v\n", err.Field, err.Err)
	}
	return result
}

func Validate(v interface{}) error {
	vValue := reflect.ValueOf(v)
	vType := vValue.Type()

	if vValue.Kind() != reflect.Struct {
		return UnsupportedInputType
	}

	validator, err := ParseRules(vType)
	if err != nil {
		return err
	}

	return validator.Validate(vValue)
}

type Validator struct {
	errors      ValidationErrors
	intRules    map[string][]IntRule
	stringRules map[string][]StringRule
}

func ParseRules(t reflect.Type) (Validator, error) {
	var result = Validator{
		make(ValidationErrors, 0),
		make(map[string][]IntRule, 0),
		make(map[string][]StringRule, 0),
	}
	var err error

	for _, field := range reflect.VisibleFields(t) {
		fieldTag, found := field.Tag.Lookup(validationTag)
		if !found {
			continue
		}

		kind := field.Type.Kind()
		if kind == reflect.Slice {
			kind = field.Type.Elem().Kind()
		}

		switch kind {
		case reflect.Int:
			result.intRules[field.Name], err = ParseIntRules(fieldTag)
		case reflect.String:
			result.stringRules[field.Name], err = ParseStringRules(fieldTag)
		default:
			return result, fmt.Errorf("unsupported type %q", kind)
		}
		if err != nil {
			return result, fmt.Errorf("field %q has invalid tag %w", field.Name, err)
		}

	}
	return result, nil
}

func (v *Validator) Validate(value reflect.Value) error {
	v.errors = make(ValidationErrors, 0)
	for fieldName, rule := range v.intRules {
		field := value.FieldByName(fieldName)
		v.checkIntField(fieldName, field, rule)
	}
	for fieldName, rule := range v.stringRules {
		field := value.FieldByName(fieldName)
		v.checkStringField(fieldName, field, rule)
	}

	return v.errors
}

func (v *Validator) checkIntField(fieldName string, field reflect.Value, rules []IntRule) {
	if !field.CanInt() {
		v.errors = append(v.errors, ValidationError{fieldName, fmt.Errorf("wrong type %T", field.Kind())})
		return
	}
	value := field.Int()
	for _, rule := range rules {
		err := rule(value)
		if err != nil {
			v.errors = append(v.errors, ValidationError{fieldName, err})
		}
	}
}

func (v *Validator) checkStringField(fieldName string, field reflect.Value, rules []StringRule) {
	value := field.String()
	for _, rule := range rules {
		err := rule(value)
		if err != nil {
			v.errors = append(v.errors, ValidationError{fieldName, err})
		}
	}
}
