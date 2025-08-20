package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GetErrorValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s is not a valid email", fe.Field())
	case "max":
		return fmt.Sprintf("%s must be less than %s", fe.Field(), fe.Param())
	case "min":
		return fmt.Sprintf("%s must be more than %s", fe.Field(), fe.Param())
	case "number":
		return fmt.Sprintf("%s must be a number", fe.Field())
	case "chickenCategory":
		return fmt.Sprintf("%s is not a valid chicken category", fe.Field())
	case "requestItemStatus":
		return fmt.Sprintf("%s is not a valid request item status", fe.Field())
	case "paymentMethod":
		return fmt.Sprintf("%s is not a valid payment method", fe.Field())
	case "itemCategory":
		return fmt.Sprintf("%s is not a valid item category", fe.Field())
	case "paymentType":
		return fmt.Sprintf("%s is not a valid payment type", fe.Field())
	case "saleUnit":
		return fmt.Sprintf("%s is not a valid sale unit", fe.Field())
	case "chickenHealthItemType":
		return fmt.Sprintf("%s is not a valid chicken health item type", fe.Field())
	case "presenceStatus":
		return fmt.Sprintf("%s is not a valid presence status", fe.Field())
	case "salaryInterval":
		return fmt.Sprintf("%s is not a valid salary interval", fe.Field())
	case "customerType":
		return fmt.Sprintf("%s is not a valid customer type", fe.Field())
	case "phoneNumber":
		return fmt.Sprintf("%s is not a valid phone number", fe.Field())
	case "eggType":
		return fmt.Sprintf("%s is not a valid egg type", fe.Field())
	case "ovenCondition":
		return fmt.Sprintf("%s is not a valid oven condition", fe.Field())
	case "feedType":
		return fmt.Sprintf("%s is not a valid feed type", fe.Field())
	case "supplierType":
		return fmt.Sprintf("%s is not a valid supplier type", fe.Field())
	case "incomeCategory":
		return fmt.Sprintf("%s is not a valid income category", fe.Field())
	case "expenseCategory":
		return fmt.Sprintf("%s is not a valid expense category", fe.Field())
	default:
		return fmt.Sprintf("%s is not valid", fe.Field())
	}
}
