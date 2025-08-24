package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type CreateItemRequest struct {
	Name          string   `json:"name" validate:"required"`
	Unit          string   `json:"unit" validate:"required"`
	Category      string   `json:"category" validate:"required,itemCategory"`
	DailySpending *float64 `json:"dailySpending"`
}

type UpdateItemRequest struct {
	Name          string   `json:"name" validate:"required"`
	Unit          string   `json:"unit" validate:"required"`
	Category      string   `json:"category" validate:"required,itemCategory"`
	DailySpending *float64 `json:"dailySpending"`
}

type GetItemFilter struct {
	Categories []param.ItemCategoryParam `query:"category"`
	ItemNames  []string                  `query:"itemNames"`
}

type ItemResponse struct {
	Id            uint64   `json:"id"`
	Name          string   `json:"name"`
	Unit          string   `json:"unit"`
	Category      string   `json:"category"`
	DailySpending *float64 `json:"dailySpending"`
}

type CreateItemPriceRequest struct {
	Category string `json:"category" validate:"required"`
	ItemId   uint64 `json:"itemId" validate:"required,number"`
	Price    string `json:"price" validate:"required"`
}

type UpdateItemPriceRequest struct {
	Category string `json:"category" validate:"required"`
	ItemId   uint64 `json:"itemId" validate:"required,number"`
	Price    string `json:"price" validate:"required"`
}

type CreateItemPriceDiscountRequest struct {
	Name                   string  `json:"name" validate:"required"`
	MinimumTransactionUser uint64  `json:"minimumTransactionUser" validate:"number,min=0"`
	TotalDiscount          float64 `json:"totalDiscount" validate:"number,gte=0,lte=100"`
}

type UpdateItemPriceDiscountRequest struct {
	Name                   string  `json:"name" validate:"required"`
	MinimumTransactionUser uint64  `json:"minimumTransactionUser" validate:"number,min=0"`
	TotalDiscount          float64 `json:"totalDiscount" validate:"number,gte=0,lte=100"`
}

type ItemPriceResponse struct {
	Id       uint64       `json:"id"`
	Category string       `json:"category"`
	Item     ItemResponse `json:"item"`
	Price    string       `json:"price"`
}

type ItemPriceDiscountResponse struct {
	Id                     uint64  `json:"id"`
	Name                   string  `json:"name"`
	MinimumTransactionUser uint64  `json:"minimumTransactionUser"`
	TotalDiscount          float64 `json:"totalDiscount"`
}
