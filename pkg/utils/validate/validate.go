package validate

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Struct(s any) error {
	return validate.Struct(s)
}
