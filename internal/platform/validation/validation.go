package validation

import (
	"github.com/go-playground/validator/v10"
)

var v = validator.New(validator.WithRequiredStructEnabled())

func CheckStruct(input any) error {
	return v.Struct(input)
}
