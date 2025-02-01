package mocks

import (
	reflect "reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// MockFieldError is a mock implementation of validator.FieldError
type MockFieldError struct {
	FieldName string
	TagName   string
}

var _ validator.FieldError = (*MockFieldError)(nil)

func (m MockFieldError) Field() string           { return m.FieldName }
func (m MockFieldError) Tag() string             { return m.TagName }
func (m MockFieldError) ActualTag() string       { return m.TagName }
func (m MockFieldError) Namespace() string       { return "" }
func (m MockFieldError) StructNamespace() string { return "" }
func (m MockFieldError) StructField() string     { return m.FieldName }
func (m MockFieldError) Param() string           { return "" }
func (m MockFieldError) Kind() reflect.Kind      { return reflect.String }
func (m MockFieldError) Type() reflect.Type      { return nil }
func (m MockFieldError) Value() interface{}      { return nil }
func (m MockFieldError) Error() string {
	return "validation failed on field " + m.FieldName + " with tag " + m.TagName
}
func (m MockFieldError) Translate(_ ut.Translator) string {
	return m.Error()
}
