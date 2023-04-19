package utils

import "astro/dto"

import "github.com/go-playground/validator/v10"

func ValidateStruct(object interface{}) []*dto.ErrorValidator {

	var validate = validator.New()

	var errors []*dto.ErrorValidator
	err := validate.Struct(object)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element dto.ErrorValidator
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
