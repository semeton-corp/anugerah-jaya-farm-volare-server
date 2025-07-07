package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type CreateStoreRequest struct {
	Name       string `json:"name"`
	LocationId uint64 `json:"locationId"`
}

type UpdateStoreRequest struct {
	Name       string `json:"name"`
	LocationId uint64 `json:"locationId"`
}

type GetStoreFilter struct {
	LocationId uint64 `query:"locationId"`
}

type StoreResponse struct {
	Id            uint64           `json:"id"`
	Name          string           `json:"name"`
	Location      LocationResponse `json:"location"`
	TotalEmployee uint64           `json:"totalEmployee"`
}

type StoreDetailResponse struct {
	Id       uint64           `json:"id"`
	Name     string           `json:"name"`
	Location LocationResponse `json:"location"`
	Users    []UserResponse   `json:"users"`
}

type CreateStoreRequestItemRequest struct {
	ItemId      uint64  `json:"itemID" validate:"required,number"`
	WarehouseId uint64  `json:"warehouseId" validate:"required,number"`
	Quantity    float64 `json:"quantity" validate:"required,number"` // ikat
}

type UpdateStoreRequestItemByWarehouseRequest struct {
	Status string `json:"status" validate:"required,requestItemStatus,oneof=Dikirim Ditolak"`
}

type UpdateStoreRequestItemByStoreRequest struct {
	Status   string  `json:"status" validate:"required,requestItemStatus,oneof=Diterima"`
	Quantity float64 `json:"quantity" validate:"required,number"`
}

type UpdateStoreRequestItemRequest struct {
	Status   string  `json:"status" validate:"required,requestItemStatus"`
	Quantity float64 `json:"quantity" validate:"required,number"`
}

type StoreRequestItemResponse struct {
	Id            uint64            `json:"id"`
	Warehouse     WarehouseResponse `json:"warehouse"`
	WarehouseItem ItemResponse      `json:"warehouseItem"`
	Store         StoreResponse     `json:"store"`
	Quantity      float64           `json:"quantity"`
	Status        string            `json:"status"`
	RequestDate   string            `json:"requestDate"`
}

type StoreRequestItemListPaginationResponse struct {
	TotalPage         uint64                     `json:"totalPage"`
	TotalData         uint64                     `json:"totalData"`
	StoreRequestItems []StoreRequestItemResponse `json:"storeRequestItems"`
}

type GetStoreRequestItemFilter struct {
	Date    param.DateParam `query:"date"`
	Page    uint64          `query:"page"`
	StoreId uint64
}

type StoreItemResponse struct {
	Store       StoreResponse `json:"store"`
	Item        ItemResponse  `json:"item"`
	Quantity    float64       `json:"quantity"`
	Description string        `json:"description"`
}

type UpdateStoreItemRequest struct {
	Quantity float64 `json:"quantity" validate:"required"`
}

type StoreItemOverview struct {
	EggStoreItemSummaries []EggStoreItemSummary `json:"eggStoreItemSummaries"`
	StoreItems            []StoreItemResponse   `json:"storeItems"`
}

type EggStoreItemSummary struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type GetStoreSaleFilter struct {
	Date          param.DateParam          `query:"date"`
	PaymentMethod param.PaymentMethodParam `query:"paymentMethod"`
	Page          uint64                   `query:"page"`
}

type StoreSaleResponse struct {
	Id               uint64                     `json:"id"`
	SendDate         string                     `json:"sentDate"`
	Customer         string                     `json:"customer"`
	Phone            string                     `json:"phone"`
	WarehouseItem    ItemResponse               `json:"warehouseItem"`
	Store            StoreResponse              `json:"store"`
	Quantity         uint64                     `json:"quantity"`
	SaleUnit         string                     `json:"saleUnit"`
	PaymentType      string                     `json:"paymentType"`
	PaymentStatus    string                     `json:"paymentStatus"`
	Price            string                     `json:"price"`
	TotalPrice       string                     `json:"totalPrice"`
	IsSend           bool                       `json:"isSend"`
	Payments         []StoreSalePaymentResponse `json:"payments"`
	RemainingPayment string                     `json:"remainingPayment"`
}

type StoreSaleListPaginationResponse struct {
	TotalPage  uint64                  `json:"totalPage"`
	TotalData  uint64                  `json:"totalData"`
	StoreSales []StoreSaleListResponse `json:"storeSales"`
}

type CreateStoreSaleRequest struct {
	Customer         string                        `json:"customer" validate:"required"`
	Phone            string                        `json:"phone" validate:"required"`
	WarehouseItemId  uint64                        `json:"warehouseItemId" validate:"required,number"`
	StoreId          uint64                        `json:"storeId" validate:"required,number"`
	Quantity         uint64                        `json:"quantity" validate:"required,number"`
	SaleUnit         string                        `json:"saleUnit" validate:"required,saleUnit"`
	Price            string                        `json:"price" validate:"required,number"`
	SendDate         string                        `json:"sendDate" validate:"required"`
	PaymentType      string                        `json:"paymentType" validate:"required,paymentType"`
	StoreSalePayment CreateStoreSalePaymentRequest `json:"storeSalePayment" validate:"required"`
}

type UpdateStoreSaleRequest struct {
	Customer        string `json:"customer" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
	WarehouseItemId uint64 `json:"warehouseItemId" validate:"required,number"`
	StoreId         uint64 `json:"storeId" validate:"required,number"`
	Quantity        uint64 `json:"quantity" validate:"required,number"`
	SaleUnit        string `json:"saleUnit" validate:"required,saleUnit"`
	Price           string `json:"price" validate:"required,number"`
	SendDate        string `json:"sendDate" validate:"required"`
	PaymentType     string `json:"paymentType" validate:"required,paymentType"`
}

type CreateStoreSalePaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateStoreSalePaymentRequest struct {
	PaymentDate  string `json:"paymentDate" validate:"required"`
	Nominal      string `json:"nominal" validate:"required,number"`
	PaymentProof string `json:"paymentProof" validate:"required,url"`
}

type StoreSaleListResponse struct {
	Id            uint64        `json:"id"`
	SendDate      string        `json:"sentDate"`
	Customer      string        `json:"customer"`
	Phone         string        `json:"phone"`
	WarehouseItem ItemResponse  `json:"warehouseItem"`
	Store         StoreResponse `json:"store"`
	Quantity      uint64        `json:"quantity"`
	SaleUnit      string        `json:"saleUnit"`
	PaymentType   string        `json:"paymentType"`
	PaymentStatus string        `json:"paymentStatus"`
	IsSend        bool          `json:"isSend"`
}

type StoreSalePaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type GetStoreItemFilter struct {
	StoreId  uint64                           `query:"storeId"`
	Category param.WarehouseItemCategoryParam `query:"category"`
}
