package validator

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("chicken_category", ValidateChickenCategory)
	validate.RegisterValidation("requestItemStatus", ValidationRequestItemStatus)
	validate.RegisterValidation("paymentMethod", ValidationPaymentMethod)
	validate.RegisterValidation("warehouseItemCategory", ValidationWarehouseItemCategory)
	validate.RegisterValidation("paymentType", ValidationPaymentType)
	validate.RegisterValidation("saleUnit", ValidationSaleUnit)

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
	case "Tunai", "Non Tunai":
		return true
	default:
		return false
	}
}

func ValidationWarehouseItemCategory(fl validator.FieldLevel) bool {
	warehouseItemCategory := fl.Field().String()
	switch warehouseItemCategory {
	case "Pakan", "Barang", "Telur", "Bahan Baku":
		return true
	default:
		return false
	}
}

func ValidationPaymentType(fl validator.FieldLevel) bool {
	paymentType := fl.Field().String()
	switch paymentType {
	case "Penuh", "Cicil":
		return true
	default:
		return false
	}
}

func ValidationSaleUnit(fl validator.FieldLevel) bool {
	saleUnit := fl.Field().String()
	switch saleUnit {
	case "Butir", "Ikat", "Karpet":
		return true
	default:
		return false
	}
}
