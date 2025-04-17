package validator

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("chicken_category", ValidateChickenCategory)
	validate.RegisterValidation("unit", ValidationUnit)
	return validate
}

func ValidateChickenCategory(fl validator.FieldLevel) bool {
	chickenCategory := fl.Field().String()
	switch chickenCategory {
	case "DOC", "Grower", "Pre Layer", "Layer", "Afkir":
		return true
	default:
		return false
	}
}

func ValidationUnit(fl validator.FieldLevel) bool {
	unit := fl.Field().String()
	switch unit {
	case "kg", "liter":
		return true
	default:
		return false
	}
}
