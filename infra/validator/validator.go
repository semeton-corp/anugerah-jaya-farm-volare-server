package validator

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("chicken_category", ValidateChickenCategory)
	return validate
}

func ValidateChickenCategory(fl validator.FieldLevel) bool {
	chickenCategory := fl.Field().String()
	switch chickenCategory {
	case "doc", "grower", "pre_layer", "layer", "afkir":
		return true
	default:
		return false
	}
}
