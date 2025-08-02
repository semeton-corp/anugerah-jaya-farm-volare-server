package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type CreateWarehouseRequest struct {
	Name         string  `json:"name" validate:"required"`
	CornCapacity float64 `json:"cornCapacity" validate:"required"`
	LocationId   uint64  `json:"locationId" validate:"required"`
}

type UpdateWarehouseRequest struct {
	Name         string  `json:"name" validate:"required"`
	CornCapacity float64 `json:"cornCapacity" validate:"required"`
	LocationId   uint64  `json:"locationId" validate:"required"`
}

type GetWarehouseFilter struct {
	LocationId uint64 `query:"locationId"`
}

type GetWarehouseItemFilter struct {
	WarehouseId uint64                  `query:"warehouseId"`
	Category    param.ItemCategoryParam `query:"category"`
	ItemNames   []string                `query:"itemNames"`
	Units       []string                `query:"units"`
}

type WarehouseResponse struct {
	Id            uint64           `json:"id"`
	Name          string           `json:"name"`
	Location      LocationResponse `json:"location"`
	CornCapacity  float64          `json:"cornCapacity"`
	TotalEmployee uint64           `json:"totalEmployee"`
}

type WarehouseDetailResponse struct {
	Id       uint64           `json:"id"`
	Name     string           `json:"name"`
	Location LocationResponse `json:"location"`
	Users    []UserResponse   `json:"users"`
}

type CreateWarehouseItemRequest struct {
	WarehouseId     uint64  `json:"warehouseId" validate:"required"`
	ItemId          uint64  `json:"itemId" validate:"required"`
	Quantity        float64 `json:"quantity" validate:"required"`
	RunOutCountDown *uint64 `json:"runOutCountDown"`
}

type UpdateWarehouseItemRequest struct {
	Quantity        float64 `json:"quantity" validate:"required"`
	RunOutCountDown *uint64 `json:"runOutCountDown"`
}

type WarehouseItemResponse struct {
	Warehouse        WarehouseResponse `json:"warehouse"`
	Item             ItemResponse      `json:"item"`
	Quantity         float64           `json:"quantity"`
	EstimationRunOut string            `json:"estimationRunOut"`
	Description      string            `json:"description"`
}

type CreateWarehouseItemProcurementRequest struct {
	WarehouseId     uint64  `json:"warehouseId" validate:"required"`
	WarehouseItemId uint64  `json:"warehouseItemId" validate:"required"`
	SupplierId      uint64  `json:"supplierId" validate:"required"`
	Quantity        float64 `json:"quantity" validate:"required"`
}

type WarehouseItempProcurementResponse struct {
	Id        uint64               `json:"id"`
	Warehouse WarehouseResponse    `json:"warehouse"`
	Item      ItemResponse         `json:"item"`
	Supplier  SupplierListResponse `json:"supplier"`
	TakenBy   string               `json:"takenBy"`
	TakenAt   string               `json:"takenAt"`
	IsTaken   bool                 `json:"isTaken"`
	Quantity  float64              `json:"quantity"`
}

type GetWarehouseItemProcurementFilter struct {
	Date    param.DateParam `query:"date"`
	IsTaken bool
}

type WarehouseOverview struct {
	TotalSafeStock    uint64                  `json:"totalSafeStock"`
	TotalDangerStock  uint64                  `json:"totalDangerStock"`
	TotalStoreRequest uint64                  `json:"totalStoreRequest"`
	EggStocks         []WarehouseItemResponse `json:"eggStocks"`
	EquipmentStocks   []WarehouseItemResponse `json:"equipmentStocks"`
}

type WarehouseItemHistoryListResponse struct {
	Id          uint64       `json:"id"`
	Item        ItemResponse `json:"item"`
	Source      string       `json:"source"`
	Destination string       `json:"destination"`
	Quantity    float64      `json:"quantity"`
	Status      string       `json:"status"`
	Time        string       `json:"time"`
}

type WarehouseItemHistoryListPaginationResponse struct {
	TotalPage              uint64                             `json:"totalPage,omitempty"`
	TotalData              uint64                             `json:"totalData,omitempty"`
	WarehouseItemHistories []WarehouseItemHistoryListResponse `json:"warehouseItemHistories"`
}

type WarehouseItemHistoryResponse struct {
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

type GetWarehouseItemHistoryFilter struct {
	Date param.DateParam `query:"date"`
	Page uint64          `query:"page"`
}

type EggWarehouseItemSummary struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type GetEggWarehouseItemSummary struct {
	WarehouseId uint64 `query:"warehouseId"`
}

type GetWarehouseSaleFilter struct {
	Date          param.DateParam          `query:"date"`
	PaymentMethod param.PaymentMethodParam `query:"paymentMethod"`
	Page          uint64                   `query:"page"`
}

type WarehouseSaleResponse struct {
	Id               uint64                         `json:"id"`
	SendDate         string                         `json:"sentDate"`
	Customer         CustomerResponse               `json:"customer"`
	WarehouseItem    ItemResponse                   `json:"item"`
	Warehouse        WarehouseResponse              `json:"warehouse"`
	Quantity         float64                        `json:"quantity"`
	SaleUnit         string                         `json:"saleUnit"`
	PaymentType      string                         `json:"paymentType"`
	PaymentStatus    string                         `json:"paymentStatus"`
	Price            string                         `json:"price"`
	TotalPrice       string                         `json:"totalPrice"`
	IsSend           bool                           `json:"isSend"`
	Payments         []WarehouseSalePaymentResponse `json:"payments"`
	RemainingPayment string                         `json:"remainingPayment"`
}

type WarehouseSaleListPaginationResponse struct {
	TotalPage      uint64                      `json:"totalPage,omitempty"`
	TotalData      uint64                      `json:"totalData,omitempty"`
	WarehouseSales []WarehouseSaleListResponse `json:"warehouseSales"`
}

type CreateWarehouseSaleRequest struct {
	CustomerId           uint64                            `json:"customerId"`
	CustomerName         string                            `json:"customerName"`
	CustomerPhoneNumber  string                            `json:"customerPhoneNumber" validate:"phoneNumber"`
	CustomerType         string                            `json:"customerType" validate:"required,customerType"`
	ItemId               uint64                            `json:"itemId" validate:"required,number"`
	WarehouseId          uint64                            `json:"warehouseId" validate:"required,number"`
	Quantity             float64                           `json:"quantity" validate:"required,number"`
	SaleUnit             string                            `json:"saleUnit" validate:"required,saleUnit"`
	Price                string                            `json:"price" validate:"required,number"`
	Discount             float64                           `json:"discount" validate:"min=0"`
	SendDate             string                            `json:"sendDate" validate:"required"`
	PaymentType          string                            `json:"paymentType" validate:"required,paymentType"`
	WarehouseSalePayment CreateWarehouseSalePaymentRequest `json:"warehouseSalePayment" validate:"required"`
}

type UpdateWarehouseSaleRequest struct {
	Quantity float64 `json:"quantity" validate:"required,number"`
	SendDate string  `json:"sendDate" validate:"required"`
	Price    string  `json:"price" validate:"required"`
	Discount float64 `json:"discount" validate:"required"`
	SaleUnit string  `json:"saleUnit" validate:"required,saleUnit"`
}

type CreateWarehouseSalePaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateWarehouseSalePaymentRequest struct {
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
}

type WarehouseSaleListResponse struct {
	Id            uint64            `json:"id"`
	OrderDate     string            `json:"orderDate"`
	SendDate      string            `json:"sentDate"`
	Customer      CustomerResponse  `json:"customer"`
	Item          ItemResponse      `json:"item"`
	Warehouse     WarehouseResponse `json:"Warehouse"`
	Quantity      float64           `json:"quantity"`
	SaleUnit      string            `json:"saleUnit"`
	PaymentStatus string            `json:"paymentStatus"`
	IsSend        bool              `json:"isSend"`
}

type WarehouseSalePaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type CreateWarehouseSaleQueueRequest struct {
	CustomerId          uint64  `json:"customerId"`
	CustomerName        string  `json:"customerName"`
	CustomerPhoneNumber string  `json:"customerPhoneNumber"`
	CustomerType        string  `json:"customerType" validate:"required,customerType"`
	ItemId              uint64  `json:"itemId" validate:"required,number"`
	WarehouseId         uint64  `json:"warehouseId" validate:"required,number"`
	Quantity            float64 `json:"quantity" validate:"required,number"`
	SaleUnit            string  `json:"saleUnit" validate:"required,saleUnit"`
}

type WarehouseSaleQueueResponse struct {
	OrderPriority uint64            `json:"OrderPriority"`
	Id            uint64            `json:"id"`
	Quantity      uint64            `json:"quantity"`
	Item          ItemResponse      `json:"item"`
	Warehouse     WarehouseResponse `json:"warehouse"`
	SaleUnit      string            `json:"saleUnit"`
	Customer      CustomerResponse  `json:"customer"`
}

type GetWarehouseSaleQueueFilter struct {
	WarehouseId uint64 `query:"warehouseId"`
}
