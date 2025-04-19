package validator

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("chicken_category", ValidateChickenCategory)
	validate.RegisterValidation("unit", ValidationUnit)
	validate.RegisterValidation("requestItemStatus", ValidationRequestItemStatus)
	validate.RegisterValidation("paymentMethod", ValidationPaymentMethod)

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

func ValidationRequestItemStatus(fl validator.FieldLevel) bool {
	requestItemStatus := fl.Field().String()
	switch requestItemStatus {
	case "Diterima", "Menunggu", "Ditolak", "Dikirim":
		return true
	default:
		return false
	}
}

func ValidationPaymentMethod(fl validator.FieldLevel) bool {
	paymentMethod := fl.Field().String()
	switch paymentMethod {
	case "Penuh", "Cicil":
		return true
	default:
		return false
	}
}
