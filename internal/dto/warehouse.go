package dto

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
	"github.com/shopspring/decimal"
)

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
	LocationId   uint64                  `query:"locationId"`
	WarehouseIds []uint64                `query:"warehouseIds"`
	WarehouseId  uint64                  `query:"warehouseId"`
	Category     param.ItemCategoryParam `query:"category"`
	ItemNames    []string                `query:"itemNames"`
	Units        []string                `query:"units"`
}

type WarehouseResponse struct {
	Id           uint64           `json:"id"`
	Name         string           `json:"name"`
	Location     LocationResponse `json:"location"`
	CornCapacity float64          `json:"cornCapacity"`
}

type WarehouseDetailResponse struct {
	Id            uint64           `json:"id"`
	Name          string           `json:"name"`
	Location      LocationResponse `json:"location"`
	CornCapacity  float64          `json:"cornCapacity"`
	IsItemsEmpty  bool             `json:"isItemsEmpty"`
	TotalEmployee uint64           `json:"totalEmployee"`
}

type WarehouseWithUsersResponse struct {
	Id            uint64             `json:"id"`
	Name          string             `json:"name"`
	Location      LocationResponse   `json:"location"`
	CornCapacity  float64            `json:"cornCapacity"`
	IsItemsEmpty  bool               `json:"isItemsEmpty"`
	TotalEmployee uint64             `json:"totalEmployee"`
	Users         []UserListResponse `json:"users"`
}

type CreateWarehouseItemRequest struct {
	WarehouseId uint64  `json:"warehouseId" validate:"required"`
	ItemId      uint64  `json:"itemId" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required"`
}

type UpdateWarehouseItemRequest struct {
	Quantity float64 `json:"quantity" validate:"min=0"`
}

type UpdateWarehouseItemCornRequest struct {
	Quantity float64 `json:"quantity" validate:"required"`
}

type WarehouseItemResponse struct {
	Warehouse        WarehouseResponse `json:"warehouse"`
	Item             ItemResponse      `json:"item"`
	Quantity         float64           `json:"quantity"`
	EstimationRunOut string            `json:"estimationRunOut"`
	Description      string            `json:"description"`
	ExpiredAt        string            `json:"expiredAt"`
}

type CreateWarehouseItemProcurementRequest struct {
	WarehouseId           uint64                                         `json:"warehouseId" validate:"required"`
	ItemId                uint64                                         `json:"itemId" validate:"required"`
	SupplierId            uint64                                         `json:"supplierId" validate:"required"`
	DailySpending         float64                                        `json:"dailySpending" validate:"required"`
	DaysNeed              uint64                                         `json:"daysNeed" validate:"required"`
	Price                 string                                         `json:"price" validate:"required"`
	EstimationArrivalDate string                                         `json:"estimationArrivalDate" validate:"required"`
	ExpiredAt             *string                                        `json:"expiredAt"`
	DeadlinePaymentDate   *string                                        `json:"deadlinePaymentDate"`
	PaymentType           string                                         `json:"paymentType" validate:"required,paymentType"`
	Payments              []CreateWarehouseItemProcurementPaymentRequest `json:"payments" validate:"dive"`
}

type CreateWarehouseItemProcurementPaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateWarehouseItemProcurementPaymentRequest struct {
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
}

type WarehouseItemProcurementListResponse struct {
	Id                            uint64               `json:"id"`
	OrderDate                     string               `json:"orderDate"`
	Warehouse                     WarehouseResponse    `json:"warehouse"`
	Item                          ItemResponse         `json:"item"`
	Supplier                      SupplierListResponse `json:"supplier"`
	IsArrived                     bool                 `json:"IsArrived"`
	Quantity                      float64              `json:"quantity"`
	EstimationArrivalDate         string               `json:"estimationArrivalDate"`
	ProcurementStatus             string               `json:"procurementStatus"`
	DeadlinePaymentDate           string               `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool                 `json:"isMoreThanDeadlinePaymentDate"`
	ExpiredAt                     string               `json:"expiredAt"`
	PaymentStatus                 string               `json:"paymentStatus"`
	PaidDate                      string               `json:"paidDate"`
}

type WarehouseItemProcurementResponse struct {
	Id                            uint64                                    `json:"id"`
	OrderDate                     string                                    `json:"orderDate"`
	Warehouse                     WarehouseResponse                         `json:"warehouse"`
	Item                          ItemResponse                              `json:"item"`
	Supplier                      SupplierListResponse                      `json:"supplier"`
	IsArrived                     bool                                      `json:"IsArrived"`
	Quantity                      float64                                   `json:"quantity"`
	RemainingPayment              string                                    `json:"remainingPayment"`
	EstimationArrivalDate         string                                    `json:"estimationArrivalDate"`
	ProcurementStatus             string                                    `json:"procurementStatus"`
	DeadlinePaymentDate           string                                    `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool                                      `json:"isMoreThanDeadlinePaymentDate"`
	PaidDate                      string                                    `json:"paidDate"`
	Price                         string                                    `json:"price"`
	DaysNeed                      uint64                                    `json:"daysNeed"`
	TotalPrice                    string                                    `json:"totalPrice"`
	ExpiredAt                     string                                    `json:"expiredAt"`
	PaymentStatus                 string                                    `json:"paymentStatus"`
	PaymentType                   string                                    `json:"paymentType"`
	Payments                      []WarehouseItemProcurementPaymentResponse `json:"payments"`
	ReceiveQuantity               *float64                                  `json:"receiveQuantity"`
	Note                          string                                    `json:"note"`
}

type WarehouseItemProcurementPaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type ArrivalConfirmationWarehouseItemProcurementRequest struct {
	Quantity float64 `json:"quantity" validate:"required"`
	Note     string  `json:"note"`
}

type WarehouseItemProcurementListPaginationResponse struct {
	TotalData                  uint64                                 `json:"totalData,omitempty"`
	TotalPage                  uint64                                 `json:"totalPage,omitempty"`
	WarehouseItemProcurementes []WarehouseItemProcurementListResponse `json:"warehouseItemProcurements"`
}

type GetWarehouseItemProcurementFilter struct {
	DeadlinePaymentStartDate param.DateParam              `query:"deadlinePaymentStartDate"`
	DeadlinePaymentEndDate   param.DateParam              `query:"deadlinePaymentEndDate"`
	PaymentStatus            param.PaymentStatusParam     `query:"paymentStatus"`
	Status                   param.ProcurementStatusParam `query:"status"`
	PaymentStatuses          []param.PaymentStatusParam   `query:"paymentStatuses"`
	Page                     uint64                       `query:"page"`
	WarehouseId              uint64                       `query:"warehouseId"`
	LocationId               uint64                       `query:"locationId"`
	ProcurementStatus        param.ProcurementStatusParam `query:"procurementStatus"`
}

type WarehouseOverview struct {
	Warehouse        WarehouseResponse           `json:"warehouse"`
	TotalSafeStock   uint64                      `json:"totalSafeStock"`
	TotalDangerStock uint64                      `json:"totalDangerStock"`
	TotalItemInOrder uint64                      `json:"totalItemInOrder"`
	EggStocks        []WarehouseItemResponse     `json:"eggStocks"`
	CornStocks       []WarehouseItemCornResponse `json:"cornStocks"`
	EquipmentStocks  []WarehouseItemResponse     `json:"equipmentStocks"`
}

type WarehouseItemHistoryListResponse struct {
	Id          uint64  `json:"id"`
	ItemName    string  `json:"itemName"`
	ItemUnit    string  `json:"itemUnit"`
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	Quantity    float64 `json:"quantity"`
	Status      string  `json:"status"`
	Time        string  `json:"time"`
}

type WarehouseItemHistoryListPaginationResponse struct {
	TotalPage              uint64                             `json:"totalPage,omitempty"`
	TotalData              uint64                             `json:"totalData,omitempty"`
	WarehouseItemHistories []WarehouseItemHistoryListResponse `json:"warehouseItemHistories"`
}

type WarehouseItemHistoryResponse struct {
	Id             uint64  `json:"id"`
	ItemName       string  `json:"itemName"`
	ItemUnit       string  `json:"itemUnit"`
	Source         string  `json:"source"`
	Destination    string  `json:"destination"`
	QuantityBefore float64 `json:"quantityBefore"`
	QuantityAfter  float64 `json:"quantityAfter"`
	Status         string  `json:"status"`
	UpdatedBy      string  `json:"updatedBy"`
	Time           string  `json:"time"`
	Date           string  `json:"date"`
}

type GetWarehouseItemHistoryFilter struct {
	Date param.DateParam `query:"date"`
	Page uint64          `query:"page"`
}

type EggWarehouseItemSummaryResponse struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type CornWarehouseItemSummaryResponse struct {
	Warehouse WarehouseResponse `json:"warehouse"`
	Name      string            `json:"name"`
	Quantity  float64           `json:"quantity"`
	Unit      string            `json:"unit"`
}

type GetEggWarehouseItemSummary struct {
	WarehouseId uint64 `query:"warehouseId"`
}

type GetWarehouseSaleFilter struct {
	DeadlinePaymentStartDate param.DateParam            `query:"deadlinePaymentStartDate"`
	DeadlinePaymentEndDate   param.DateParam            `query:"deadlinePaymentEndDate"`
	Date                     param.DateParam            `query:"date"`
	PaymentStatus            param.PaymentStatusParam   `query:"paymentStatus"`
	Page                     uint64                     `query:"page"`
	StartDate                param.DateParam            `query:"startDate"`
	EndDate                  param.DateParam            `query:"endDate"`
	WarehouseId              uint64                     `query:"warehouseId"`
	LocationId               uint64                     `query:"locationId"`
	ItemId                   uint64                     `query:"itemId"`
	PaymentStatuses          []param.PaymentStatusParam `query:"paymentStatuses"`
}

type WarehouseSaleResponse struct {
	Id                            uint64                         `json:"id"`
	SendDate                      string                         `json:"sentDate"`
	Customer                      CustomerResponse               `json:"customer"`
	WarehouseItem                 ItemResponse                   `json:"item"`
	Warehouse                     WarehouseResponse              `json:"warehouse"`
	Quantity                      float64                        `json:"quantity"`
	SaleUnit                      string                         `json:"saleUnit"`
	PaymentType                   string                         `json:"paymentType"`
	PaymentStatus                 string                         `json:"paymentStatus"`
	Price                         string                         `json:"price"`
	TotalPrice                    string                         `json:"totalPrice"`
	IsSend                        bool                           `json:"isSend"`
	Payments                      []WarehouseSalePaymentResponse `json:"payments"`
	RemainingPayment              string                         `json:"remainingPayment"`
	DeadlinePaymentDate           string                         `json:"deadlinePaymentDate"`
	PaidDate                      string                         `json:"paidDate"`
	IsMoreThanDeadlinePaymentDate bool                           `json:"isMoreThanDeadlinePaymentDate"`
}

type WarehouseSaleListPaginationResponse struct {
	TotalPage      uint64                      `json:"totalPage,omitempty"`
	TotalData      uint64                      `json:"totalData,omitempty"`
	WarehouseSales []WarehouseSaleListResponse `json:"warehouseSales"`
}

type CreateWarehouseSaleRequest struct {
	CustomerId          uint64                              `json:"customerId"`
	CustomerName        string                              `json:"customerName"`
	CustomerPhoneNumber string                              `json:"customerPhoneNumber" validate:"phoneNumber"`
	CustomerType        string                              `json:"customerType" validate:"required,customerType"`
	ItemId              uint64                              `json:"itemId" validate:"required,number"`
	WarehouseId         uint64                              `json:"warehouseId" validate:"required,number"`
	Quantity            float64                             `json:"quantity" validate:"required,number"`
	SaleUnit            string                              `json:"saleUnit" validate:"required,saleUnit"`
	Price               string                              `json:"price" validate:"required,number"`
	Discount            float64                             `json:"discount" validate:"min=0"`
	SendDate            string                              `json:"sendDate" validate:"required"`
	PaymentType         string                              `json:"paymentType" validate:"required,paymentType"`
	Payments            []CreateWarehouseSalePaymentRequest `json:"payments" validate:"dive"`
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
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateWarehouseSalePaymentRequest struct {
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
}

type WarehouseSaleListResponse struct {
	Id                            uint64            `json:"id"`
	OrderDate                     string            `json:"orderDate"`
	SendDate                      string            `json:"sentDate"`
	Customer                      CustomerResponse  `json:"customer"`
	Item                          ItemResponse      `json:"item"`
	Warehouse                     WarehouseResponse `json:"warehouse"`
	Quantity                      float64           `json:"quantity"`
	SaleUnit                      string            `json:"saleUnit"`
	PaymentStatus                 string            `json:"paymentStatus"`
	IsSend                        bool              `json:"isSend"`
	CreatedAt                     time.Time         `json:"-"`
	TotalPrice                    decimal.Decimal   `json:"-"`
	DeadlinePaymentDate           string            `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool              `json:"isMoreThanDeadlinePaymentDate"`
	PaidDate                      string            `json:"paidDate"`
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
	OrderPriority   uint64            `json:"orderPriority"`
	Id              uint64            `json:"id"`
	Quantity        float64           `json:"quantity"`
	Item            ItemResponse      `json:"item"`
	Warehouse       WarehouseResponse `json:"warehouse"`
	SaleUnit        string            `json:"saleUnit"`
	Customer        CustomerResponse  `json:"customer"`
	TotalAllocation float64           `json:"totalAllocation"`
}

type GetWarehouseSaleQueueFilter struct {
	WarehouseId uint64 `query:"warehouseId"`
}

type CreateWarehouseItemProcurementDraftRequest struct {
	WarehouseId   uint64  `json:"warehouseId" validate:"required"`
	ItemId        uint64  `json:"itemId" validate:"required"`
	SupplierId    uint64  `json:"supplierId" validate:"required"`
	DailySpending float64 `json:"dailySpending" validate:"required"`
	DaysNeed      uint64  `json:"daysNeed" validate:"required"`
	Price         string  `json:"price" validate:"required"`
}

type UpdateWarehouseItemProcurementDraftRequest struct {
	WarehouseId   uint64  `json:"warehouseId" validate:"required"`
	ItemId        uint64  `json:"itemId" validate:"required"`
	SupplierId    *uint64 `json:"supplierId"`
	DailySpending float64 `json:"dailySpending" validate:"required"`
	DaysNeed      uint64  `json:"daysNeed" validate:"required"`
	Price         string  `json:"price" validate:"required"`
}

type WarehouseItemProcurementDraftResponse struct {
	Id            uint64               `json:"id"`
	InputDate     string               `json:"inputDate"`
	Warehouse     WarehouseResponse    `json:"warehouse"`
	Item          ItemResponse         `json:"item"`
	Supplier      SupplierListResponse `json:"supplier"`
	DailySpending float64              `json:"dailySpending"`
	DaysNeed      uint64               `json:"daysNeed"`
	Quantity      float64              `json:"quantity"`
	Price         string               `json:"price"`
	TotalPrice    string               `json:"totalPrice"`
}

type CreateWarehouseItemCornProcurementDraftRequest struct {
	WarehouseId               uint64  `json:"warehouseId" validate:"required"`
	SupplierId                uint64  `json:"supplierId" validate:"required"`
	OvenCondition             *string `json:"ovenCondition"`
	CornWaterLevel            float64 `json:"cornWaterLevel" validate:"required"`
	IsOvenCanOperateInNearDay *bool   `json:"isOvenCanOperateInNearDay"`
	Quantity                  float64 `json:"quantity" validate:"required"`
	Price                     string  `json:"price" validate:"required"`
	Discount                  float64 `json:"discount" validate:"min=0"`
}

type UpdateWarehouseItemCornProcurementDraftRequest struct {
	WarehouseId               uint64  `json:"warehouseId" validate:"required"`
	SupplierId                uint64  `json:"supplierId" validate:"required"`
	OvenCondition             *string `json:"ovenCondition"`
	CornWaterLevel            float64 `json:"cornWaterLevel" validate:"required"`
	IsOvenCanOperateInNearDay *bool   `json:"isOvenCanOperateInNearDay"`
	Quantity                  float64 `json:"quantity" validate:"required"`
	Price                     string  `json:"price" validate:"required"`
	Discount                  float64 `json:"discount" validate:"min=0"`
}

type WarehouseItemCornProcurementDraftResponse struct {
	Id                        uint64               `json:"id"`
	InputDate                 string               `json:"inputDate"`
	Warehouse                 WarehouseResponse    `json:"warehouse"`
	Supplier                  SupplierListResponse `json:"supplier"`
	Item                      ItemResponse         `json:"item"`
	OvenCondition             string               `json:"ovenCondition"`
	CornWaterLevel            *float64             `json:"cornWaterLevel"`
	IsOvenCanOperateInNearDay *bool                `json:"isOvenCanOperateInNearDay"`
	Quantity                  float64              `json:"quantity"`
	Discount                  *float64             `json:"discount"`
	Price                     string               `json:"price"`
	TotalPrice                string               `json:"totalPrice"`
}

type CreateWarehouseItemCornProcurementRequest struct {
	WarehouseId               uint64                                             `json:"warehouseId" validate:"required"`
	SupplierId                uint64                                             `json:"supplierId" validate:"required"`
	OvenCondition             *string                                            `json:"ovenCondition"`
	CornWaterLevel            float64                                            `json:"cornWaterLevel" validate:"required"`
	IsOvenCanOperateInNearDay *bool                                              `json:"isOvenCanOperateInNearDay"`
	Quantity                  float64                                            `json:"quantity" validate:"required"`
	Price                     string                                             `json:"price" validate:"required"`
	ExpiredAt                 string                                             `json:"expiredAt" validate:"required"`
	Discount                  float64                                            `json:"discount" validate:"min=0"`
	DeadlinePaymentDate       *string                                            `json:"deadlinePaymentDate"`
	PaymentType               string                                             `json:"paymentType" validate:"required,paymentType"`
	Payments                  []CreateWarehouseItemCornProcurementPaymentRequest `json:"payments" validate:"dive"`
}

type WarehouseItemCornProcurementListResponse struct {
	Id                            uint64               `json:"id"`
	OrderDate                     string               `json:"orderDate"`
	Warehouse                     WarehouseResponse    `json:"warehouse"`
	Supplier                      SupplierListResponse `json:"supplier"`
	Item                          ItemResponse         `json:"item"`
	TotalPrice                    string               `json:"totalPrice"`
	ProcurementStatus             string               `json:"procurementStatus"`
	IsArrived                     bool                 `json:"IsArrived"`
	PaymentStatus                 string               `json:"paymentStatus"`
	DeadlinePaymentDate           string               `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool                 `json:"isMoreThanDeadlinePaymentDate"`
	Quantity                      float64              `json:"quantity"`
	Discount                      float64              `json:"discount"`
	Price                         string               `json:"price"`
	PaidDate                      string               `json:"paidDate"`
	ReceiveQuantity               *float64             `json:"recieveQuantity"`
}

type WarehouseItemCornProcurementListPaginationResponse struct {
	TotalData                     uint64                                     `json:"totalData,omitempty"`
	TotalPage                     uint64                                     `json:"totalPage,omitempty"`
	WarehouseItemCornProcurements []WarehouseItemCornProcurementListResponse `json:"warehouseItemCornProcurements"`
}

type WarehouseItemCornProcurementResponse struct {
	Id                            uint64                                        `json:"id"`
	Warehouse                     WarehouseResponse                             `json:"warehouse"`
	Supplier                      SupplierListResponse                          `json:"supplier"`
	Item                          ItemResponse                                  `json:"item"`
	IsArrived                     bool                                          `json:"IsArrived"`
	OvenCondition                 string                                        `json:"ovenCondition"`
	CornWaterLevel                float64                                       `json:"cornWaterLevel"`
	ProcurementStatus             string                                        `json:"procurementStatus"`
	IsOvenCanOperateInNearDay     *bool                                         `json:"isOvenCanOperateInNearDay"`
	Price                         string                                        `json:"price" validate:"required"`
	Quantity                      float64                                       `json:"quantity"`
	TotalPrice                    string                                        `json:"totalPrice"`
	RemainingPayment              string                                        `json:"remainingPayment"`
	DeadlinePaymentDate           string                                        `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool                                          `json:"isMoreThanDeadlinePaymentDate"`
	PaidDate                      string                                        `json:"paidDate"`
	Discount                      float64                                       `json:"discount"`
	PaymentStatus                 string                                        `json:"paymentStatus"`
	PaymentType                   string                                        `json:"paymentType"`
	Payments                      []WarehouseItemCornProcurementPaymentResponse `json:"payments"`
	ReceieveQuantity              *float64                                      `json:"recieveQuantity"`
	ExpiredAt                     string                                        `json:"expiredAt"`
	Date                          string                                        `json:"date"`
}

type WarehouseItemCornResponse struct {
	Id        uint64               `json:"id"`
	OrderDate string               `json:"orderDate"`
	Quantity  float64              `json:"quantity"`
	Item      ItemResponse         `json:"item"`
	Supplier  SupplierListResponse `json:"supplier"`
	ExpiredAt string               `json:"expiredAt"`
}

type GetWarehouseItemCornFilter struct {
	WarehouseId      uint64 `query:"warehouseId"`
	WithZeroQuantity *bool  `query:"WithZeroQuantity"`
	FromNewest       bool   `query:"fromNewest"`
}

type GetWarehouseItemCornProcurementFilter struct {
	DeadlinePaymentStartDate param.DateParam              `query:"deadlinePaymentStartDate"`
	DeadlinePaymentEndDate   param.DateParam              `query:"deadlinePaymentEndDate"`
	PaymentStatus            param.PaymentStatusParam     `query:"paymentStatus"`
	PaymentStatuses          []param.PaymentStatusParam   `query:"paymentStatuses"`
	Page                     uint64                       `query:"page"`
	WarehouseId              uint64                       `query:"warehouseId"`
	LocationId               uint64                       `query:"locationId"`
	ProcurementStatus        param.ProcurementStatusParam `query:"procurementStatus"`
}

type CreateWarehouseItemCornProcurementPaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateWarehouseItemCornProcurementPaymentRequest struct {
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
}

type WarehouseItemCornProcurementPaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type ArrivalConfirmationWarehouseItemCornProcurementRequest struct {
	Quantity float64 `json:"quantity" validate:"required"`
	Note     string  `json:"note"`
}

type WarehouseItemCornPriceResponse struct {
	Id          uint64  `json:"id"`
	BottomLimit float64 `json:"bottomLimit"`
	UpperLimit  float64 `json:"upperLimit"`
	Discount    float64 `json:"discount"`
}

type CreateRawFeedRequest struct {
	WarehouseId  uint64                     `json:"warehouseId" validate:"required"`
	CornQuantity float64                    `json:"cornQuantity" validate:"required"`
	CornPrice    string                     `json:"cornPrice" validate:"required"`
	DaysNeed     uint64                     `json:"daysNeed" validate:"required"`
	RawMaterials []CreateRawMaterialRequest `json:"rawMaterials" validate:"dive"`
}

type CreateRawMaterialRequest struct {
	ItemId        uint64  `json:"itemId" validate:"required"`
	Quantity      float64 `json:"quantity" validate:"required"`
	Price         string  `json:"price" validate:"required"`
	DailySpending float64 `json:"dailySpending" validate:"required"`
}

type CreateReadyToEatFeedRequest struct {
	WarehouseId   uint64  `json:"warehouseId" validate:"required"`
	ItemId        uint64  `json:"itemId" validate:"required"`
	DaysNeed      uint64  `json:"daysNeed" validate:"required"`
	Price         string  `json:"price" validate:"required"`
	DailySpending float64 `json:"dailySpending" validate:"required"`
}

type ReduceFeedRequest struct {
	ItemId       uint64  `json:"itemId"`
	ItemCategory string  `json:"itemCategory"`
	Quantity     float64 `json:"quantity"`
}

type GetWarehouseSalePaymentFilter struct {
	StartDate  param.DateParam `query:"startDate"`
	EndDate    param.DateParam `query:"endDate"`
	LocationId uint64          `query:"locationId"`
	Date       param.DateParam `query:"date"`
}

type GetStoreSalePaymentFilter struct {
	StartDate  param.DateParam `query:"startDate"`
	Page       uint64          `query:"page"`
	EndDate    param.DateParam `query:"endDate"`
	LocationId uint64          `query:"locationId"`
	StoreId    uint64          `query:"storeId"`
	Date       param.DateParam `query:"date"`
}

type GetAfkirChickenSalePaymentFilter struct {
	StartDate  param.DateParam `query:"startDate"`
	EndDate    param.DateParam `query:"endDate"`
	LocationId uint64          `query:"locationId"`
	Date       param.DateParam `query:"date"`
}

type GetWarehouseItemProcurementPaymentFilter struct {
	StartDate  param.DateParam `query:"startDate"`
	EndDate    param.DateParam `query:"endDate"`
	LocationId uint64          `query:"locationId"`
	Date       param.DateParam `query:"date"`
}

type GetWarehouseItemCornProcurementPaymentFilter struct {
	StartDate  param.DateParam `query:"startDate"`
	EndDate    param.DateParam `query:"endDate"`
	LocationId uint64          `query:"locationId"`
	Date       param.DateParam `query:"date"`
}

type GetWarehouseItemProcurementDraftFilter struct {
	WarehouseId  uint64                  `query:"warehouseId"`
	Date         param.DateParam         `query:"date"`
	ItemCategory param.ItemCategoryParam `query:"itemCategory"`
}

type GetWarehouseItemCornProcurementDraftFilter struct {
	WarehouseId uint64          `query:"warehouseId"`
	Date        param.DateParam `query:"date"`
}
