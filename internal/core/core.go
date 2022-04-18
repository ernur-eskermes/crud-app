package core

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func validateStruct(st interface{}) []*ErrorResponse {
	var res []*ErrorResponse

	if err := validate.Struct(st); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, err := range validationErrors {
				var element ErrorResponse
				element.FailedField = err.StructNamespace()
				element.Tag = err.Tag()
				element.Value = err.Param()
				res = append(res, &element)
			}
		}
	}

	return res
}
