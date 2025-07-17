package validator

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("chickenCategory", ValidateChickenCategory)
	validate.RegisterValidation("requestItemStatus", ValidationRequestItemStatus)
	validate.RegisterValidation("paymentMethod", ValidationPaymentMethod)
	validate.RegisterValidation("itemCategory", ValidationItemCategory)
	validate.RegisterValidation("paymentType", ValidationPaymentType)
	validate.RegisterValidation("saleUnit", ValidationSaleUnit)
	validate.RegisterValidation("chickenHealthItemType", ValidationChickenHealthItemType)
	validate.RegisterValidation("presenceStatus", ValidationPresenceStatus)
	validate.RegisterValidation("salaryInterval", ValidationSalaryInterval)
	validate.RegisterValidation("customerType", ValidationCustomerType)
	validate.RegisterValidation("phoneNumber", ValidationPhoneNumber)

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

func ValidationItemCategory(fl validator.FieldLevel) bool {
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
	case "Kg", "Ikat", "Plastik":
		return true
	default:
		return false
	}
}

func ValidationChickenHealthItemType(fl validator.FieldLevel) bool {
	chickenHealthItemType := fl.Field().String()
	switch chickenHealthItemType {
	case "Obat", "Vaksin Kondisional", "Vaksin Rutin":
		return true
	default:
		return false
	}
}

func ValidationPresenceStatus(fl validator.FieldLevel) bool {
	presenceStatus := fl.Field().String()
	switch presenceStatus {
	case "Hadir", "Sakit", "Izin":
		return true
	default:
		return false
	}
}

func ValidationSalaryInterval(fl validator.FieldLevel) bool {
	presenceStatus := fl.Field().String()
	switch presenceStatus {
	case "Harian", "Bulanan":
		return true
	default:
		return false
	}
}

func ValidationCustomerType(fl validator.FieldLevel) bool {
	presenceStatus := fl.Field().String()
	switch presenceStatus {
	case "Pelanggan Baru", "Pelanggan Lama":
		return true
	default:
		return false
	}
}

func ValidationPhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()

	return phoneNumber[:2] == "08"
}
