package model

import "github.com/go-playground/validator/v10"

func RegisterValidators(validate *validator.Validate) {
	validate.RegisterValidation("state", EmptyOrBelongsToEnum(States))
}

func EmptyOrBelongsToEnum(enum []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		if value == "" {
			return true
		}

		for _, s := range enum {
			if s == value {
				return true
			}
		}

		return false
	}
}
