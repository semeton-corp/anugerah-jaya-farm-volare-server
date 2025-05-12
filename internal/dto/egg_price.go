package dto

type CreateEggPriceRequest struct {
	Category        string `json:"category" validate:"required"`
	WarehouseItemId uint64 `json:"WarehouseItemId" validate:"required,number"`
	Price           string `json:"price" validate:"required"`
}

type UpdateEggPriceRequest struct {
	Category        string `json:"category" validate:"required"`
	WarehouseItemId uint64 `json:"WarehouseItemId" validate:"required,number"`
	Price           string `json:"price" validate:"required"`
}

type CreateEggPriceDiscountRequest struct {
	Name                   string  `json:"name" validate:"required"`
	MinimumTransactionUser uint64  `json:"minimumTransactionUser" validate:"required,number"`
	TotalDiscount          float64 `json:"totalDiscount" validate:"required"`
}

type UpdateEggPriceDiscountRequest struct {
	Name                   string  `json:"name" validate:"required"`
	MinimumTransactionUser uint64  `json:"minimumTransactionUser" validate:"required,number"`
	TotalDiscount          float64 `json:"totalDiscount" validate:"required"`
}

type EggPriceResponse struct {
	Id            uint64                `json:"id"`
	Category      string                `json:"category"`
	WarehouseItem WarehouseItemResponse `json:"warehouse"`
	Price         string                `json:"price"`
}

type EggPriceDiscountResponse struct {
	Id                     uint64  `json:"id"`
	Name                   string  `json:"name"`
	MinimumTransactionUser uint64  `json:"minimumTransactionUser"`
	TotalDiscount          float64 `json:"totalDiscount"`
}
