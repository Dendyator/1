package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "нет ошибок валидации"
	}
	var sb strings.Builder
	for _, e := range v {
		sb.WriteString(e.Field + ": " + e.Err.Error() + "\n")
	}
	return sb.String()
}

func Validate(v interface{}) error {
	valErrors := ValidationErrors{}
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return errors.New("ожидается структура")
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		fieldErrors := validateField(fieldValue, validateTag, field.Name)
		valErrors = append(valErrors, fieldErrors...)
	}

	if len(valErrors) > 0 {
		return valErrors
	}
	return nil
}

func validateField(value reflect.Value, tag, fieldName string) ValidationErrors {
	var errs ValidationErrors
	errorTracker := map[string]bool{}

	switch value.Kind() { //nolint:exhaustive
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			elemErrors := validateField(value.Index(i), tag, fieldName)
			for _, err := range elemErrors {
				if !errorTracker[err.Err.Error()] {
					errorTracker[err.Err.Error()] = true
					errs = append(errs, err)
				}
			}
		}
	case reflect.String, reflect.Int:
		tagParts := strings.Split(tag, "|")
		for _, part := range tagParts {
			if err := validateRule(value, strings.TrimSpace(part)); err != nil {
				if !errorTracker[err.Error()] {
					errorTracker[err.Error()] = true
					errs = append(errs, ValidationError{Field: fieldName, Err: err})
				}
			}
		}
	default:
	}
	return errs
}

func validateRule(value reflect.Value, rule string) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		return validateString(value.String(), rule)
	case reflect.Int:
		return validateInt(int(value.Int()), rule)
	default:
		return nil
	}
}

func validateString(val, rule string) error {
	switch {
	case strings.HasPrefix(rule, "len:"):
		expectedLen, _ := strconv.Atoi(rule[4:])
		if len(val) != expectedLen {
			return errors.New("длина должна быть " + strconv.Itoa(expectedLen))
		}
	case strings.HasPrefix(rule, "regexp:"):
		re, err := regexp.Compile(rule[7:])
		if err != nil {
			return err
		}
		if !re.MatchString(val) {
			return errors.New("не соответствует регулярному выражению")
		}
	case strings.HasPrefix(rule, "in:"):
		options := strings.Split(rule[3:], ",")
		for _, opt := range options {
			if val == opt {
				return nil
			}
		}
		return errors.New("значение не входит в варианты")
	}
	return nil
}

func validateInt(val int, rule string) error {
	switch {
	case strings.HasPrefix(rule, "min:"):
		minimum, _ := strconv.Atoi(rule[4:])
		if val < minimum {
			return errors.New("не может быть меньше " + strconv.Itoa(minimum))
		}
	case strings.HasPrefix(rule, "max:"):
		maximum, _ := strconv.Atoi(rule[4:])
		if val > maximum {
			return errors.New("не может быть больше " + strconv.Itoa(maximum))
		}
	case strings.HasPrefix(rule, "in:"):
		options := strings.Split(rule[3:], ",")
		for _, opt := range options {
			intOpt, _ := strconv.Atoi(opt)
			if val == intOpt {
				return nil
			}
		}
		return errors.New("значение не входит в варианты")
	}
	return nil
}
