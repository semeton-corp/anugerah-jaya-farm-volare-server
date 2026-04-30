package dto

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type CreateChickenMonitoringRequest struct {
	ChickenCageId     uint64  `json:"chickenCageId" validate:"required"`
	TotalSickChicken  uint64  `json:"totalSickChicken" validate:"number,min=0"`
	TotalDeathChicken uint64  `json:"totalDeathChicken" validate:"number,min=0"`
	TotalFeed         float64 `json:"totalFeed" validate:"number,min=0"`
	Note              string  `json:"note"`
}

type UpdateChickenMonitoringRequest struct {
	ChickenCageId     uint64  `json:"chickenCageId" validate:"required"`
	TotalSickChicken  uint64  `json:"totalSickChicken" validate:"number,min=0"`
	TotalDeathChicken uint64  `json:"totalDeathChicken" validate:"number,min=0"`
	TotalFeed         float64 `json:"totalFeed" validate:"number,min=0"`
	Note              string  `json:"note"`
}

type ChickenMonitoringResponse struct {
	Id                 uint64              `json:"id"`
	ChickenCage        ChickenCageResponse `json:"chickenCage"`
	TotalLiveChicken   uint64              `json:"totalLiveChicken"`
	TotalSickChicken   uint64              `json:"totalSickChicken"`
	TotalDeatchChicken uint64              `json:"totalDeathChicken"`
	TotalFeed          float64             `json:"totalFeed"`
	Note               string              `json:"note"`
}

type ChickenMonitoringListResponse struct {
	Id                uint64              `json:"id"`
	ChickenCage       ChickenCageResponse `json:"chickenCage"`
	TotalLiveChicken  uint64              `json:"totalLiveChicken"`
	TotalSickChicken  uint64              `json:"totalSickChicken"`
	TotalDeathChicken uint64              `json:"totalDeathChicken"`
	TotalFeed         float64             `json:"totalFeed"`
	MortalityRate     float64             `json:"mortalityRate"`
}

type CreateChickenHealthItemRequest struct {
	Name       string  `json:"name" validate:"required"`
	Type       string  `json:"type" validate:"required,chickenHealthItemType"`
	ChickenAge *uint64 `json:"chickenAge"`
	Note       string  `json:"note"`
}

type UpdateChickenHealthItemRequest struct {
	Name       string  `json:"name" validate:"required"`
	Type       string  `json:"type" validate:"required,chickenHealthItemType"`
	ChickenAge *uint64 `json:"chickenAge"`
	Note       string  `json:"note"`
}

type ChickenHealthItemResponse struct {
	Id              uint64  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	ChickenAge      *uint64 `json:"chickenAge"`
	ChickenCategory *string `json:"chickenCategory"`
	Note            string  `json:"note"`
}

type GetChickenHealthItemFilter struct {
	Type param.ChickenHealthItemTypeParam `query:"type"`
}

type CreateChickenHealthMonitoringRequest struct {
	ChickenCageId  uint64  `json:"chickenCageId" validate:"required"`
	HealthItemName string  `json:"healthItemName" validate:"required"`
	Type           string  `json:"type" validate:"required,chickenHealthItemType"`
	Dose           float64 `json:"dose" validate:"required"`
	Unit           string  `json:"unit" validate:"required"`
	Disease        *string `json:"disease"`
}

type UpdateChickenHealthMonitoringRequest struct {
	ChickenCageId  uint64  `json:"chickenCageId" validate:"required"`
	HealthItemName string  `json:"healthItemName" validate:"required"`
	Type           string  `json:"type" validate:"required,chickenHealthItemType"`
	Dose           float64 `json:"dose" validate:"required"`
	Unit           string  `json:"unit" validate:"required"`
	Disease        *string `json:"disease"`
}

type ChickenHealthMonitoringResponse struct {
	Id              uint64  `json:"id"`
	HealthItemName  string  `json:"healthItemName"`
	Type            string  `json:"type"`
	Dose            float64 `json:"dose"`
	Unit            string  `json:"unit"`
	Disease         string  `json:"disease"`
	Date            string  `json:"date"`
	ChickenAge      uint64  `json:"chickenAge"`
	ChickenCategory string  `json:"chickenCategory"`
	CreatedAt       string  `json:"createdAt"`
}

type ChickenHealthMonitoringDetailResponse struct {
	ChickenCage              ChickenCageResponse               `json:"chickenCage"`
	ChickenHealthMonitorings []ChickenHealthMonitoringResponse `json:"chickenHealthMonitorings"`
}

type GetChickenMonitoringFilter struct {
	Date          param.DateParam `query:"date"`
	LocationId    uint64          `query:"locationId"`
	CageId        uint64
	ChickenCageId uint64
	StartDate     param.DateParam
	EndDate       param.DateParam
}

type GetChickenOverviewFilter struct {
	LocationId        uint64                       `query:"locationId"`
	CageId            uint64                       `query:"cageId"`
	OverviewGraphTime param.OverviewGraphTimeParam `query:"overviewGraphTime"`
	Year              int                          `query:"year"`
}

type ChickenDetailOverview struct {
	TotalLiveChicken    uint64  `json:"totalLiveChicken"`
	TotalSickChicken    uint64  `json:"totalSickChicken"`
	TotalDeathChicken   uint64  `json:"totalDeathChicken"`
	TotalKPIPerformance float64 `json:"totalKPIPerformance"`
}

type ChickenGraphResponse struct {
	Key          string `json:"key"`
	SickChicken  uint64 `json:"sickChicken"`
	DeathChicken uint64 `json:"deathChicken"`
}

type ChickenBarChartResponse struct {
	ChickenDOC       float64 `json:"chickenDOC"`
	ChickenGrower    float64 `json:"chickenGrower"`
	ChickentPreLayer float64 `json:"chickenPreLayer"`
	ChickenLayer     float64 `json:"chickenLayer"`
	ChickenAfkir     float64 `json:"chickenAfkir"`
}

type ChickenOverviewResponse struct {
	ChickenDetail ChickenDetailOverview   `json:"chickenDetail"`
	ChickenGraphs []ChickenGraphResponse  `json:"chickenGraphs"`
	ChickenPie    ChickenBarChartResponse `json:"chickenPie"`
}

type CreateChickenProcurementDraftRequest struct {
	CageId     uint64 `json:"cageId" validate:"required"`
	SupplierId uint64 `json:"supplierId" validate:"required"`
	Quantity   uint64 `json:"quantity" validate:"required"`
	TotalPrice string `json:"totalPrice" validate:"required"`
}

type UpdateChickenProcurementDraftRequest struct {
	CageId     uint64 `json:"cageId" validate:"required"`
	SupplierId uint64 `json:"supplierId" validate:"required"`
	Quantity   uint64 `json:"quantity" validate:"required"`
	TotalPrice string `json:"totalPrice" validate:"required"`
}

type ChickenProcurementDraftResponse struct {
	Id         uint64           `json:"id"`
	InputDate  string           `json:"inputDate"`
	Cage       CageResponse     `json:"cage"`
	Supplier   SupplierResponse `json:"supplier"`
	Quantity   uint64           `json:"quantity"`
	TotalPrice string           `json:"totalPrice"`
}

type ConfirmationChickenProcurementRequest struct {
	Quantity            uint64                                   `json:"quantity" validate:"required"`
	PaymentType         string                                   `json:"paymentType" validate:"required,paymentType"`
	TotalPrice          string                                   `json:"totalPrice" validate:"required"`
	DeadlinePaymentDate *string                                  `json:"deadlinePaymentDate"`
	EstimateArrivalDate string                                   `json:"estimationArrivalDate" validate:"required"`
	Payments            []CreateChickenProcurementPaymentRequest `json:"payments"`
}

type ChickenProcurementResponse struct {
	Id                            uint64                              `json:"id"`
	OrderDate                     string                              `json:"orderDate"`
	Cage                          CageResponse                        `json:"cage"`
	Supplier                      SupplierListResponse                `json:"supplier"`
	Quantity                      uint64                              `json:"quantity"`
	ReceiveQuantity               *uint64                             `json:"receiveQuantity"`
	TotalPrice                    string                              `json:"totalPrice"`
	EstimationArrivalDate         string                              `json:"estimationArrivalDate"`
	Payments                      []ChickenProcurementPaymentResponse `json:"payments"`
	PaymentStatus                 string                              `json:"paymentStatus"`
	PaymentType                   string                              `json:"paymentType"`
	RemainingPayment              string                              `json:"remainingPayment"`
	IsArrived                     bool                                `json:"IsArrived"`
	DeadlinePaymentDate           string                              `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool                                `json:"isMoreThanDeadlinePaymentDate"`
	ProcurementStatus             string                              `json:"procurementStatus"`
	PaidDate                      string                              `json:"paidDate"`
	Note                          string                              `json:"note"`
}

type ChickenProcurementListResponse struct {
	Id                            uint64               `json:"id"`
	OrderDate                     string               `json:"orderDate"`
	Quantity                      uint64               `json:"quantity"`
	Cage                          CageResponse         `json:"cage"`
	Supplier                      SupplierListResponse `json:"supplier"`
	EstimationArrivalDate         string               `json:"estimationArrivalDate"`
	PaymentStatus                 string               `json:"paymentStatus"`
	IsArrived                     bool                 `json:"IsArrived"`
	PaymentType                   string               `json:"paymentType"`
	DeadlinePaymentDate           string               `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool                 `json:"isMoreThanDeadlinePaymentDate"`
	ProcurementStatus             string               `json:"procurementStatus"`
	TotalPrice                    string               `json:"totalPrice"`
	PaidDate                      string               `json:"paidDate"`
}

type ChickenProcurementListPaginationResponse struct {
	TotalPage           uint64                           `json:"totalPage,omitempty"`
	TotalData           uint64                           `json:"totalData,omitempty"`
	ChickenProcurements []ChickenProcurementListResponse `json:"chickenProcurements"`
}

type GetChickenProcurementFilter struct {
	DeadlinePaymentStartDate param.DateParam            `query:"deadlinePaymentStartDate"`
	DeadlinePaymentEndDate   param.DateParam            `query:"deadlinePaymentEndDate"`
	PaymentStatuses          []param.PaymentStatusParam `query:"paymentStatuses"`
	PaymentStatus            param.PaymentStatusParam   `query:"paymentStatus"`
	Page                     uint64                     `query:"page"`
	LocationId               uint64
}

type CreateChickenProcurementPaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateChickenProcurementPaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type ChickenProcurementPaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type ArrivalConfirmationChickenProcurementRequest struct {
	Quantity uint64 `json:"quantity" validate:"required"`
	Note     string `json:"note"`
}

// Note : AverageEggWeight only for good egg
type ChickenPerformanceSummaryResponse struct {
	FeedConsumption      float64 `json:"foodConsumption"`
	AverageEggWeight     float64 `json:"averageEggWeight"`
	AverageFCR           float64 `json:"averageFCR"`
	AverageHDP           float64 `json:"averageHDP"`
	AverageMortalityRate float64 `json:"averageMortalityRate"`
}

type ChickenCagePerformanceSummaryResponse struct {
	TotalProductiveCage uint64 `json:"totalProductiveCage"`
	TotalAfkirCage      uint64 `json:"totalAfkirCage"`
}

type WarehouseItemSummaryResponse struct {
	TotalSafeItem    uint64 `json:"totalSafeItem"`
	TotalNotSafeItem uint64 `json:"totalNotSafeItem"`
}

type IncomeAndExpenseBarChartResponse struct {
	Key     string `json:"key"`
	Income  string `json:"income"`
	Expense string `json:"expense"`
}

type ChickenAndWarehouseOverviewResponse struct {
	ChickenPerformanceSummary     ChickenPerformanceSummaryResponse     `json:"chickenPerformanceSummary"`
	ChickenCagePerformanceSummary ChickenCagePerformanceSummaryResponse `json:"chickenCagePerformanceSummary"`
	WarehouseItemSummary          WarehouseItemSummaryResponse          `json:"warehouseItemSummary"`
	ChickenBarCharts              ChickenBarChartResponse               `json:"chickenBarChart"`
	ChickenGraphs                 []ChickenGraphResponse                `json:"chickenGraphs"`
}

type ChickenAndCompanyOverviewResponse struct {
	ChickenPerformanceSummary            ChickenPerformanceSummaryResponse  `json:"chickenPerformanceSummary"`
	IncomeAndExpensePerformanceBarCharts []IncomeAndExpenseBarChartResponse `json:"incomeAndExpensePerformanceBarCharts"`
	ChickenBarCharts                     ChickenBarChartResponse            `json:"chickenBarCharts"`
	BEPGoodEgg                           float64                            `json:"bepGoodEgg"`
	MarginOfSafety                       float64                            `json:"marginOfSafety"`
	RCRatio                              float64                            `json:"rcRatio"`
}

type GetChickenAndCompanyOverviewFilter struct {
	LocationId uint64 `query:"locationId"`
	CageId     uint64 `query:"cageId"`
	Year       uint64 `query:"year" validate:"required"`
}

type GetChickenAndWarehouseOverviewFilter struct {
	LocationId        uint64                       `query:"locationId"`
	CageId            uint64                       `query:"cageId"`
	OverviewGraphTime param.OverviewGraphTimeParam `query:"overviewGraphTime" validate:"required"`
	WarehouseId       uint64                       `query:"warehouseId"`
	Year              int                          `query:"year"`
}

type GetChickenPerformanceOverviewFilter struct {
	LocationId              uint64                       `query:"locationId"`
	CageId                  uint64                       `query:"cageId"`
	LabelCompanyPerformance string                       `query:"labelCompanyPerformance"`
	OverviewGraphTime       param.OverviewGraphTimeParam `query:"overviewGraphTime" validate:"required"`
}

type CreateAfkirChickenCustomerRequest struct {
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required,phoneNumber"`
	Address     string `json:"address" validate:"required"`
}

type UpdateAfkirChickenCustomerRequest struct {
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required,phoneNumber"`
	Address     string `json:"address" validate:"required"`
}

type AfkirChickenCustomerListResponse struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	LatestPrice string `json:"latestPrice"`
}

type AfkirChickenCustomerResponse struct {
	Id                uint64                         `json:"id"`
	Name              string                         `json:"name"`
	PhoneNumber       string                         `json:"phoneNumber"`
	Address           string                         `json:"address"`
	LatestPrice       string                         `json:"latestPrice"`
	AfkirChickenSales []AfkirChickenSaleListResponse `json:"afkirChickenSales"`
}

type AfkirChickenSaleListResponse struct {
	Id                            uint64                           `json:"id"`
	SellDate                      string                           `json:"sellDate"`
	AfkirChickenCustomer          AfkirChickenCustomerListResponse `json:"afkirChickenCustomer"`
	ChickenAge                    uint64                           `json:"chickenAge"`
	TotalSellChicken              uint64                           `json:"totalSellChicken"`
	PricePerChicken               string                           `json:"pricePerChicken"`
	TotalPrice                    string                           `json:"totalPrice"`
	PaymentStatus                 string                           `json:"paymentStatus"`
	DeadlinePaymentDate           string                           `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool                             `json:"isMoreThanDeadlinePaymentDate"`
	PaidDate                      string                           `json:"paidDate"`
	TakenAt                       string                           `json:"takenAt"`
	IsTaken                       bool                             `json:"isTaken"`
}

type AfkirChickenSaleListPaginationResponse struct {
	TotalPage         uint64                         `json:"totalPage,omitempty"`
	TotalData         uint64                         `json:"totalData,omitempty"`
	AfkirChickenSales []AfkirChickenSaleListResponse `json:"afkirChickenSales"`
}

type CreateAfkirChickenSaleDraftRequest struct {
	ChickenCageId          uint64 `json:"chickenCageId" validate:"required"`
	AfkirChickenCustomerId uint64 `json:"afkirChickenCustomerId" validate:"required"`
	TotalSellChicken       uint64 `json:"totalSellChicken" validate:"required"`
	PricePerChicken        string `json:"pricePerChicken" validate:"required"`
}

type UpdateAfkirChickenSaleDraftRequest struct {
	ChickenCageId          uint64 `json:"chickenCageId" validate:"required"`
	AfkirChickenCustomerId uint64 `json:"afkirChickenCustomerId" validate:"required"`
	TotalSellChicken       uint64 `json:"totalSellChicken" validate:"required"`
	PricePerChicken        string `json:"pricePerChicken" validate:"required"`
}

type AfkirChickenSaleDraftResponse struct {
	Id                   uint64                           `json:"id"`
	InputDate            string                           `json:"inputDate"`
	ChickenCage          ChickenCageResponse              `json:"chickenCage"`
	AfkirChickenCustomer AfkirChickenCustomerListResponse `json:"afkirChickenCustomer"`
	TotalSellChicken     uint64                           `json:"totalSellChicken"`
	PricePerChicken      string                           `json:"pricePerChicken"`
	TotalPrice           string                           `json:"totalPrice"`
}

type CreateAfkirChickenSaleRequest struct {
	AfkirChickenCustomerId uint64                                 `json:"afkirChickenCustomerId" validate:"required"`
	ChickenCageId          uint64                                 `json:"chickenCageId" validate:"required"`
	TotalSellChicken       uint64                                 `json:"totalSellChicken" validate:"required"`
	PricePerChicken        string                                 `json:"pricePerChicken" validate:"required"`
	PaymentType            string                                 `json:"paymentType" validate:"required,paymentType"`
	TakenAt                string                                 `json:"takenAt" validate:"required"`
	Payments               []CreateAfkirChickenSalePaymentRequest `json:"payments"`
}

type AfkirChickenSaleResponse struct {
	Id                            uint64                            `json:"id"`
	SellDate                      string                            `json:"sellDate"`
	AfkirChickenCustomer          AfkirChickenCustomerListResponse  `json:"afkirChickenCustomer"`
	ChickenCage                   ChickenCageResponse               `json:"chickenCage"`
	TotalSellChicken              uint64                            `json:"totalSellChicken"`
	PricePerChicken               string                            `json:"pricePerChicken"`
	TotalPrice                    string                            `json:"totalPrice"`
	PaymentType                   string                            `json:"paymentType"`
	ChickenAge                    uint64                            `json:"chickenAge"`
	Payments                      []AfkirChickenSalePaymentResponse `json:"payments"`
	PaymentStatus                 string                            `json:"paymentStatus"`
	RemainingPayment              string                            `json:"remainingPayment"`
	DeadlinePaymentDate           string                            `json:"deadlinePaymentDate"`
	IsMoreThanDeadlinePaymentDate bool                              `json:"isMoreThanDeadlinePaymentDate"`
	PaidDate                      string                            `json:"paidDate"`
	TakenAt                       string                            `json:"takenAt"`
	IsTaken                       bool                              `json:"isTaken"`
}

type CreateAfkirChickenSalePaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateAfkirChickenSalePaymentRequest struct {
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
}

type AfkirChickenSalePaymentResponse struct {
	Id            uint64 `json:"id"`
	Date          string `json:"date"`
	Nominal       string `json:"nominal"`
	Remaining     string `json:"remaining"`
	PaymentMethod string `json:"paymentMethod"`
	PaymentProof  string `json:"paymentProof"`
}

type GetAfkirChickenSaleFilter struct {
	DeadlinePaymentStartDate param.DateParam            `query:"deadlinePaymentStartDate"`
	DeadlinePaymentEndDate   param.DateParam            `query:"deadlinePaymentEndDate"`
	PaymentStatus            param.PaymentStatusParam   `query:"paymentStatus"`
	Page                     uint64                     `query:"page"`
	LocationId               uint64                     `query:"locationId"`
	PaymentStatuses          []param.PaymentStatusParam `query:"paymentStatuses"`
}

type GetChickenPerformanceFilter struct {
	Date       param.DateParam `query:"date" validate:"required"`
	LocationId uint64          `query:"locationId"`
	CageId     uint64          `query:"cageId"`
}

type ChickenPerformanceResponse struct {
	CageName                     string  `json:"cageName"`
	ChickenCategory              string  `json:"chickenCategory"`
	ChickenAge                   uint64  `json:"chickenAge"`
	TotalChicken                 uint64  `json:"totalChicken"`
	TotalEgg                     uint64  `json:"totalGoodEgg"`
	AverageConsumptionPerChicken float64 `json:"averageConsumptionPerChicken"`
	AverageWeightPerEgg          float64 `json:"averageWeightPerEgg"`
	FCR                          float64 `json:"fcr"`
	HDP                          float64 `json:"hdp"`
	MortalityRate                float64 `json:"mortalityRate"`
	Productivity                 string  `json:"productivity"`
}

type GetChickenProcurementPaymentFilter struct {
	StartDate  param.DateParam `query:"startDate"`
	EndDate    param.DateParam `query:"endDate"`
	LocationId uint64          `query:"locationId"`
	Date       param.DateParam `query:"date"`
}
