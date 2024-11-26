package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in           interface{}
		name         string
		expectedErrs ValidationErrors
	}{
		{
			name: "Valid User",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    25,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "09876543211"},
			},
			expectedErrs: nil,
		},
		{
			name: "Invalid User",
			in: User{
				ID:     "short-id",
				Name:   "Jane Doe",
				Age:    17,
				Email:  "jane.doe@example",
				Role:   "user",
				Phones: []string{"1234567", "111111"},
			},
			expectedErrs: ValidationErrors{
				{Field: "ID", Err: errors.New("длина должна быть 36")},
				{Field: "Age", Err: errors.New("не может быть меньше 18")},
				{Field: "Email", Err: errors.New("не соответствует регулярному выражению")},
				{Field: "Role", Err: errors.New("значение не входит в варианты")},
				{Field: "Phones", Err: errors.New("длина должна быть 11")},
			},
		},
		{
			name:         "Invalid App Version",
			in:           App{Version: "1.0"},
			expectedErrs: ValidationErrors{{Field: "Version", Err: errors.New("длина должна быть 5")}},
		},
		{
			name:         "Invalid Response Code",
			in:           Response{Code: 201},
			expectedErrs: ValidationErrors{{Field: "Code", Err: errors.New("значение не входит в варианты")}},
		},
		{
			name:         "Valid Response",
			in:           Response{Code: 200},
			expectedErrs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			validateTestCase(t, tt.in, tt.expectedErrs)
		})
	}
}

func validateTestCase(t *testing.T, in interface{}, expectedErrs ValidationErrors) {
	t.Helper()
	err := Validate(in)

	if err == nil {
		if expectedErrs != nil {
			t.Errorf("expected error: %v, got: nil", expectedErrs)
		}
		return
	}

	var valErrs ValidationErrors
	if errors.As(err, &valErrs) {
		if !compareValidationErrors(valErrs, expectedErrs) {
			t.Errorf("expected error: %v, got: %v", expectedErrs, valErrs)
		}
		return
	}

	t.Errorf("expected ValidationErrors type, got: %v", err)
}

func compareValidationErrors(err1, err2 ValidationErrors) bool {
	if len(err1) != len(err2) {
		return false
	}
	for i := range err1 {
		if err1[i].Field != err2[i].Field || err1[i].Err.Error() != err2[i].Err.Error() {
			return false
		}
	}
	return true
}
