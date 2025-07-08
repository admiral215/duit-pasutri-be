package utils

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateStruct(data interface{}) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = msgForTag(err.Tag(), err)
	}
	return errors
}

func msgForTag(tag string, fe validator.FieldError) string {
	switch tag {
	case "required":
		return fe.Field() + " wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return fe.Field() + " minimal " + fe.Param() + " karakter"
	case "max":
		return fe.Field() + " maksimal " + fe.Param() + " karakter"
	default:
		return fe.Field() + " tidak valid"
	}
}
