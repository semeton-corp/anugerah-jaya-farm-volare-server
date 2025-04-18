package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type StoreResponse struct {
	Id       uint64           `json:"id"`
	Name     string           `json:"name"`
	Location LocationResponse `json:"location"`
}

type CreateStoreRequestItemRequest struct {
	WarehouseItemId uint64 `json:"warehouseItemId" validate:"required,number"`
	WarehouseId     uint64 `json:"warehouseId" validate:"required,number"`
	StoreId         uint64 `json:"storeId" validate:"required,number"`
	Quantity        uint64 `json:"quantity" validate:"required,number"`
}

type UpdateStoreRequestItemByWarehouseRequest struct {
	Status string `json:"status" validate:"required,requestItemStatus,oneof=Dikirim Ditolak"`
}

type UpdateStoreRequestItemByStoreRequest struct {
	Status   string `json:"status" validate:"required,requestItemStatus,oneof=Diterima"`
	Quantity uint64 `json:"quantity" validate:"required,number"`
}

type UpdateStoreRequestItemRequest struct {
	Status   string `json:"status" validate:"required,requestItemStatus"`
	Quantity uint64 `json:"quantity" validate:"required,number"`
}

type StoreRequestItemResponse struct {
	Id            uint64                `json:"id"`
	Warehouse     WarehouseResponse     `json:"warehouse"`
	WarehouseItem WarehouseItemResponse `json:"warehouseItem"`
	Store         StoreResponse         `json:"store"`
	Quantity      uint64                `json:"quantity"`
	Status        string                `json:"status"`
	RequestDate   string                `json:"requestDate"`
}

type GetStoreRequestItemFilter struct {
	Date param.DateParam `query:"date"`
	Page uint64          `query:"page"`
}

type StoreItemResponse struct {
	Store         StoreResponse         `json:"store"`
	WarehouseItem WarehouseItemResponse `json:"warehouseItem"`
	Quantity      uint64                `json:"quantity"`
	Description   string                `json:"description"`
}

type GetStoreSaleFilter struct {
	Date param.DateParam `query:"date"`
	Page uint64          `query:"page"`
}

type StoreSaleResponse struct {
	Id               uint64                     `json:"id"`
	SendDate         string                     `json:"sentDate"`
	Customer         string                     `json:"customer"`
	Phone            string                     `json:"phone"`
	WarehouseItem    WarehouseItemResponse      `json:"warehouseItem"`
	Store            StoreResponse              `json:"store"`
	Quantity         uint64                     `json:"quantity"`
	PaymentMethod    string                     `json:"paymentMethod"`
	IsSend           bool                       `json:"isSend"`
	Payments         []StoreSalePaymentResponse `json:"payments"`
	RemainingPayment string                     `json:"remainingPayment"`
}

type CreateStoreSaleRequest struct {
	Customer         string                        `json:"customer" validate:"required"`
	Phone            string                        `json:"phone" validate:"required"`
	WarehouseItemId  uint64                        `json:"warehouseItemId" validate:"required,number"`
	StoreId          uint64                        `json:"storeId" validate:"required,number"`
	Quantity         uint64                        `json:"quantity" validate:"required,number"`
	Price            string                        `json:"price" validate:"required,number"`
	SendDate         string                        `json:"sendDate" validate:"required"`
	PaymentMethod    string                        `json:"paymentMethod" validate:"required,paymentMethod"`
	StoreSalePayment CreateStoreSalePaymentRequest `json:"storeSalePayment"`
}

type CreateStoreSalePaymentRequest struct {
	Nominal      string `json:"nominal" validate:"required,number"`
	PaymentProof string `json:"paymentProof" validate:"required,url"`
}

type StoreSaleListResponse struct {
	Id            uint64                `json:"id"`
	SendDate      string                `json:"sentDate"`
	Customer      string                `json:"customer"`
	Phone         string                `json:"phone"`
	WarehouseItem WarehouseItemResponse `json:"warehouseItem"`
	Store         StoreResponse         `json:"store"`
	Quantity      uint64                `json:"quantity"`
	PaymentMethod string                `json:"paymentMethod"`
	IsSend        bool                  `json:"isSend"`
}

type StoreSalePaymentResponse struct {
	Id           uint64 `json:"id"`
	Date         string `json:"date"`
	Nominal      string `json:"nominal"`
	PaymentProof string `json:"paymentProof"`
}

type UpdateStoreSaleRequest struct {
	IsSend bool `json:"isSend"`
}
