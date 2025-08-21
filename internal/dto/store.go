package dto

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
	"github.com/shopspring/decimal"
)

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
	StoreId     uint64  `json:"storeId" validate:"required,number"`
	ItemId      uint64  `json:"itemId" validate:"required,number"`
	WarehouseId uint64  `json:"warehouseId" validate:"required,number"`
	Quantity    float64 `json:"quantity" validate:"required,number"` // ikat
}

type UpdateStoreRequestItemRequest struct {
	Status string `json:"status" validate:"required,requestItemStatus"`
}

type WarehouseConfirmationStoreRequestItem struct {
	StoreId       uint64  `json:"storeId" validate:"required"`
	Quantity      float64 `json:"quantity" validate:"required"`
	WarehouseNote string  `json:"warehouseNote"`
}

type StoreConfirmationStoreRequestItem struct {
	Quantity  float64 `json:"quantity" validate:"required"`
	StoreNote string  `json:"storeNote"`
}

type SortingStoreRequestItemRequest struct {
	BrokenEggInButir uint64  `json:"brokenEggInButir" validate:"required"`
	BrokenEggInKg    float64 `json:"brokenEggInKg" validate:"required"`
}

type WarehouseSendUnknownStoreRequestItem struct {
	StoreId uint64 `json:"storeId" validate:"required"`
}

type StoreRequestItemResponse struct {
	Id                   uint64            `json:"id"`
	Warehouse            WarehouseResponse `json:"warehouse"`
	Store                StoreResponse     `json:"store"`
	Item                 ItemResponse      `json:"item"`
	Quantity             float64           `json:"quantity"`
	Status               string            `json:"status"`
	RequestDate          string            `json:"requestDate"`
	RecieveDate          string            `json:"recieveDate"`
	IsSorted             bool              `json:"isSorted"`
	WarehouseFulFillment float64           `json:"warehoseFulFillment"`
	CreatedBy            string            `json:"createdBy,omitempty"`
}

type StoreRequestItemListPaginationResponse struct {
	TotalPage         uint64                     `json:"totalPage,omitempty"`
	TotalData         uint64                     `json:"totalData,omitempty"`
	StoreRequestItems []StoreRequestItemResponse `json:"storeRequestItems"`
}

type GetStoreRequestItemFilter struct {
	Date        param.DateParam `query:"date"`
	Page        uint64          `query:"page"`
	StoreId     uint64          `query:"storeId"`
	WarehouseId uint64          `query:"warehouseId"`
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

type StoreItemHistoryListResponse struct {
	Id          uint64       `json:"id"`
	Item        ItemResponse `json:"item"`
	Source      string       `json:"source"`
	Destination string       `json:"destination"`
	Quantity    float64      `json:"quantity"`
	Status      string       `json:"status"`
	Time        string       `json:"time"`
}

type StoreItemHistoryResponse struct {
	Id             uint64       `json:"id"`
	Item           ItemResponse `json:"item"`
	Source         string       `json:"source"`
	Destination    string       `json:"destination"`
	QuantityBefore float64      `json:"quantityBefore"`
	QuantityAfter  float64      `json:"quantityAfter"`
	Status         string       `json:"status"`
	UpdatedBy      string       `json:"updatedBy"`
	Time           string       `json:"time"`
	Date           string       `json:"date"`
}

type GetStoreItemHistoryFilter struct {
	Date param.DateParam `query:"date"`
	Page uint64          `query:"page"`
}

type StoreItemHistoryListPaginationResponse struct {
	TotalPage          uint64                         `json:"totalPage"`
	TotalData          uint64                         `json:"totalData"`
	StoreItemHistories []StoreItemHistoryListResponse `json:"storeItemHistories"`
}

type GetStoreSaleFilter struct {
	DeadlinePaymentStartDate param.DateParam          `query:"deadlinePaymentStartDate"`
	DeadlinePaymentEndDate   param.DateParam          `query:"deadlinePaymentEndDate"`
	Date                     param.DateParam          `query:"date"`
	PaymentStatus            param.PaymentStatusParam `query:"paymentStatus"`
	Page                     uint64                   `query:"page"`
	StoreId                  uint64
	StartDate                param.DateParam
	EndDate                  param.DateParam
	ItemId                   uint64
}

type GetStoreSaleQueueFilter struct {
	StoreId uint64 `query:"storeId"`
}

type StoreSaleResponse struct {
	Id                            uint64                     `json:"id"`
	SendDate                      string                     `json:"sentDate"`
	Customer                      CustomerResponse           `json:"customer"`
	WarehouseItem                 ItemResponse               `json:"item"`
	Store                         StoreResponse              `json:"store"`
	Quantity                      float64                    `json:"quantity"`
	SaleUnit                      string                     `json:"saleUnit"`
	PaymentType                   string                     `json:"paymentType"`
	PaymentStatus                 string                     `json:"paymentStatus"`
	Price                         string                     `json:"price"`
	TotalPrice                    string                     `json:"totalPrice"`
	IsSend                        bool                       `json:"isSend"`
	Payments                      []StoreSalePaymentResponse `json:"payments"`
	RemainingPayment              string                     `json:"remainingPayment"`
	DeadlinePaymentDate           string                     `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool                       `json:"isMoreThanDeadlinePaymentDate"`
}

type StoreSaleListPaginationResponse struct {
	TotalPage  uint64                  `json:"totalPage,omitempty"`
	TotalData  uint64                  `json:"totalData,omitempty"`
	StoreSales []StoreSaleListResponse `json:"storeSales"`
}

type CreateStoreSaleRequest struct {
	CustomerId          uint64                          `json:"customerId"`
	CustomerName        string                          `json:"customerName"`
	CustomerPhoneNumber string                          `json:"customerPhoneNumber"`
	CustomerType        string                          `json:"customerType" validate:"required,customerType"`
	ItemId              uint64                          `json:"itemId" validate:"required,number"`
	StoreId             uint64                          `json:"storeId" validate:"required,number"`
	Quantity            float64                         `json:"quantity" validate:"required,number"`
	SaleUnit            string                          `json:"saleUnit" validate:"required,saleUnit"`
	Price               string                          `json:"price" validate:"required,number"`
	Discount            float64                         `json:"discount" validate:"min=0"`
	SendDate            string                          `json:"sendDate" validate:"required"`
	PaymentType         string                          `json:"paymentType" validate:"required,paymentType"`
	Payments            []CreateStoreSalePaymentRequest `json:"payments" validate:"required,dive"`
}

type UpdateStoreSaleRequest struct {
	Quantity float64 `json:"quantity" validate:"required,number"`
	SendDate string  `json:"sendDate" validate:"required"`
	Price    string  `json:"price" validate:"required,number"`
	Discount float64 `json:"discount" validate:"required"`
	SaleUnit string  `json:"saleUnit" validate:"required,saleUnit"`
}

type CreateStoreSalePaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateStoreSalePaymentRequest struct {
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
}

type StoreSaleListResponse struct {
	Id                            uint64           `json:"id"`
	OrderDate                     string           `json:"orderDate"`
	SendDate                      string           `json:"sentDate"`
	Customer                      CustomerResponse `json:"customer"`
	Item                          ItemResponse     `json:"item"`
	Store                         StoreResponse    `json:"store"`
	Quantity                      float64          `json:"quantity"`
	SaleUnit                      string           `json:"saleUnit"`
	PaymentStatus                 string           `json:"paymentStatus"`
	IsSend                        bool             `json:"isSend"`
	DeadlinePaymentDate           string           `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool             `json:"isMoreThanDeadlinePaymentDate"`
	CreatedAt                     time.Time        `json:"-"`
	TotalPrice                    decimal.Decimal  `json:"-"`
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
	StoreId   uint64                  `query:"storeId"`
	Category  param.ItemCategoryParam `query:"category"`
	ItemNames []string                `query:"itemNames"`
	Units     []string                `query:"units"`
}

type StoreOverview struct {
	StoreOverviewDetail StoreOverviewDetail  `json:"storeOverviewDetail"`
	StoreGraphs         []StoreGraphResponse `json:"storeGraphs"`
}

type StoreOverviewDetail struct {
	TotalReceivables   string  `json:"totalReceivables"`
	TotalIncome        string  `json:"totalIncome"`
	GoodEggInKg        float64 `json:"goodEggInKg"`
	GoodEggInIkat      float64 `json:"goodEggInIkat"`
	CrackedEggInKg     float64 `json:"crackedEggInKg"`
	CrackedEggInIkat   float64 `json:"crackedEggInIkat"`
	BrokenEggInPlastik float64 `json:"brokenEggInPlastik"`
}

type StoreGraphResponse struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type GetStoreOverviewFilter struct {
	ItemId            uint64                       `query:"itemId" validate:"required"`
	StoreId           uint64                       `query:"storeId" validate:"required"`
	OverviewGraphTime param.OverviewGraphTimeParam `query:"overviewGraphTime" validate:"required"`
	Year              uint64                       `query:"year" validate:"required"`
	Month             param.MonthParam             `query:"month" validate:"required"`
}

type CreateStoreSaleQueueRequest struct {
	CustomerId          uint64  `json:"customerId"`
	CustomerName        string  `json:"customerName"`
	CustomerPhoneNumber string  `json:"customerPhoneNumber"`
	CustomerType        string  `json:"customerType" validate:"required,customerType"`
	ItemId              uint64  `json:"itemId" validate:"required,number"`
	StoreId             uint64  `json:"storeId" validate:"required,number"`
	Quantity            float64 `json:"quantity" validate:"required,number"`
	SaleUnit            string  `json:"saleUnit" validate:"required,saleUnit"`
}

type StoreSaleQueueResponse struct {
	OrderPriority uint64           `json:"orderPriority"`
	Id            uint64           `json:"id"`
	Quantity      float64          `json:"quantity"`
	Item          ItemResponse     `json:"item"`
	Store         StoreResponse    `json:"store"`
	SaleUnit      string           `json:"saleUnit"`
	Customer      CustomerResponse `json:"customer"`
	CustomerType  string           `json:"customerType"`
}

type StoreItemSummaryResponse struct {
	TotalSafeItem    uint64 `json:"totalSafeItem"`
	TotalNotSafeItem uint64 `json:"totalNotSafeItem"`
}
