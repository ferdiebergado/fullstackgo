//go:generate mockgen -destination=mocks/validator_mock.go -package=mocks . Validator
package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(s any) error
}

func Instance() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}
