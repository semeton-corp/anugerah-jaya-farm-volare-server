package validator

import (
	"github.com/go-playground/validator/v10"
)

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
	validate.RegisterValidation("eggType", ValidationEggType)
	validate.RegisterValidation("ovenCondition", ValidationOvenCondition)
	validate.RegisterValidation("feedType", ValidationFeedType)
	validate.RegisterValidation("supplierType", ValidationSupplierType)
	validate.RegisterValidation("incomeCategory", ValidationIncomeCategory)
	validate.RegisterValidation("expenseCategory", ValidationExpenseCategory)
	validate.RegisterValidation("receivablesCategory", ValidationReceivablesCategory)
	validate.RegisterValidation("debtCategory", ValidationDebtCategory)
	validate.RegisterValidation("approvalStatus", ValidationApprovalStatus)
	validate.RegisterValidation("locationType", ValidationLocationType)

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
	case "Diterima", "Menunggu", "Ditolak", "Dikirim", "Dibatalkan":
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
	case "Telur", "Barang", "Bahan Baku Adukan", "Bahan Baku Adukan - Jagung", "Ayam", "Pakan Jadi":
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

	if phoneNumber != "" {
		return phoneNumber[:2] == "08"
	}

	return true
}

func ValidationEggType(fl validator.FieldLevel) bool {
	eggType := fl.Field().String()
	switch eggType {
	case "Telur OK", "Telur Retak", "Telur Bonyok":
		return true
	default:
		return false
	}
}

func ValidationOvenCondition(fl validator.FieldLevel) bool {
	eggType := fl.Field().String()
	switch eggType {
	case "Hidup", "Mati":
		return true
	default:
		return false
	}
}

func ValidationFeedType(fl validator.FieldLevel) bool {
	eggType := fl.Field().String()
	switch eggType {
	case "Pakan Jadi", "Pakan Adukan":
		return true
	default:
		return false
	}
}

func ValidationSupplierType(fl validator.FieldLevel) bool {
	supplierType := fl.Field().String()
	switch supplierType {
	case "Barang", "Ayam DOC":
		return true
	default:
		return false
	}
}

func ValidationIncomeCategory(fl validator.FieldLevel) bool {
	incomeCategory := fl.Field().String()
	switch incomeCategory {
	case "Penjualan Ayam Afkir", "Penjualan Telur Toko", "Penjualan Telur Gudang", "Semua":
		return true
	default:
		return false
	}
}

func ValidationExpenseCategory(fl validator.FieldLevel) bool {
	expenseCategory := fl.Field().String()
	switch expenseCategory {
	case "Operasional", "Lain-lain", "Semua", "Pegawai", "Pengadaan Ayam DOC", "Pengadaan Barang", "Pengadaan Jagung":
		return true
	default:
		return false
	}
}

func ValidationReceivablesCategory(fl validator.FieldLevel) bool {
	expenseCategory := fl.Field().String()
	switch expenseCategory {
	case "Penjualan Ayam Afkir", "Penjualan Telur Toko", "Penjualan Telur Gudang", "Semua", "Kasbon":
		return true
	default:
		return false
	}
}

func ValidationDebtCategory(fl validator.FieldLevel) bool {
	expenseCategory := fl.Field().String()
	switch expenseCategory {
	case "Pengadaan Ayam DOC", "Pengadaan Gudang", "Pengadaan Jagung", "Semua":
		return true
	default:
		return false
	}
}

func ValidationApprovalStatus(fl validator.FieldLevel) bool {
	approvalStatus := fl.Field().String()
	switch approvalStatus {
	case "Ditolak", "Disetujui":
		return true
	default:
		return false
	}
}

func ValidationLocationType(fl validator.FieldLevel) bool {
	locationType := fl.Field().String()
	switch locationType {
	case "Kandang", "Gudang", "Toko", "Site":
		return true
	default:
		return false
	}
}
