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
	Date       param.DateParam `query:"date"`
	LocationId uint64          `query:"locationId"`
	CageId     uint64
	StartDate  param.DateParam
	EndDate    param.DateParam
}

type GetChickenOverviewFilter struct {
	LocationId        uint64                       `query:"locationId"`
	CageId            uint64                       `query:"cageId"`
	OverviewGraphTime param.OverviewGraphTimeParam `query:"overviewGraphTime"`
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
	ChickentPreLayer float64 `json:"chickentPreLayer"`
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
	Price      string `json:"price" validate:"required"`
}

type UpdateChickenProcurementDraftRequest struct {
	CageId     uint64 `json:"cageId" validate:"required"`
	SupplierId uint64 `json:"supplierId" validate:"required"`
	Quantity   uint64 `json:"quantity" validate:"required"`
	Price      string `json:"price" validate:"required"`
}

type ChickenProcurementDraftResponse struct {
	Id         uint64           `json:"id"`
	InputDate  string           `json:"inputDate"`
	Cage       CageResponse     `json:"cage"`
	Supplier   SupplierResponse `json:"supplier"`
	Quantity   uint64           `json:"quantity"`
	Price      string           `json:"price"`
	TotalPrice string           `json:"totalPrice"`
}

type ConfirmationChickenProcurementRequest struct {
	Quantity            uint64                                   `json:"quantity"`
	Price               string                                   `json:"price"`
	EstimateArrivalDate string                                   `json:"estimationArrivalDate"`
	Payments            []CreateChickenProcurementPaymentRequest `json:"payments"`
}

type ChickenProcurementResponse struct {
	Id                    uint64                              `json:"id"`
	OrderDate             string                              `json:"orderDate"`
	Cage                  CageResponse                        `json:"cage"`
	Supplier              SupplierListResponse                `json:"supplier"`
	Quantity              uint64                              `json:"quantity"`
	TotalPrice            string                              `json:"totalPrice"`
	EstimationArrivalDate string                              `json:"estimationArrivalDate"`
	Payments              []ChickenProcurementPaymentResponse `json:"payments"`
	PaymentStatus         string                              `json:"paymentStatus"`
	RemainingPayment      string                              `json:"remainingPayment"`
}

type ChickenProcurementListResponse struct {
	Id                    uint64               `json:"id"`
	OrderDate             string               `json:"orderDate"`
	Quantity              uint64               `json:"quantity"`
	Supplier              SupplierListResponse `json:"supplier"`
	EstimationArrivalDate string               `json:"estimationArrivalDate"`
	PaymentStatus         string               `json:"paymentStatus"`
	IsArrived             bool                 `json:"IsArrived"`
}

type ChickenProcurementListPaginationResponse struct {
	TotalPage           uint64                           `json:"totalPage,omitempty"`
	TotalData           uint64                           `json:"totalData,omitempty"`
	ChickenProcurements []ChickenProcurementListResponse `json:"chickenProcurements"`
}

type GetChickenProcurementFilter struct {
	PaymentStatus param.PaymentStatusParam `query:"paymentStatus"`
	Page          uint64                   `query:"page"`
}

type CreateChickenProcurementPaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateChickenProcurementPaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
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

type ChickenPerformanceSummaryResponse struct {
	FeedConsumption  float64 `json:"foodConsumption"`
	AverageEggWeight float64 `json:"averageEggWeight"`
	AverageFCS       float64 `json:"averageFCS"`
	AverageHDP       float64 `json:"averageHDP"`
	AverageMortality float64 `json:"averageMortality"`
}

type ChickenCagePerformanceSummaryResponse struct {
	TotalProductiveCage uint64 `json:"totalProductiveCage"`
	TotalCheckCage      uint64 `json:"totalCheckCage"`
	TotalNotSafeCage    uint64 `json:"totalNotSafeCage"`
}

type WarehouseItemSummaryResponse struct {
	TotalSafeItem    uint64 `json:"totalSafeItem"`
	TotalNotSafeItem uint64 `json:"totalNotSafeItem"`
	TotalSentOffItem uint64 `json:"totalSentOffItem"`
}

type CompanyPerformanceBarChartResponse struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type ChickenPerformanceOverviewResponse struct {
	ChickenPerformanceSummary     ChickenPerformanceSummaryResponse     `json:"chickenPerformanceSummary"`
	ChickenBarCharts              ChickenBarChartResponse               `json:"chickenBarCharts"`
	CompanyPerformanceBarCharts   CompanyPerformanceBarChartResponse    `json:"companyPerformanceBarCharts"`
	ChickenCagePerformanceSummary ChickenCagePerformanceSummaryResponse `json:"chickenCagePerformanceSummary"`
	WarehouseItemSummary          WarehouseItemSummaryResponse          `json:"warehouseItemSummary"`
	ChickenGraphs                 ChickenGraphResponse                  `json:"chickenGraphs"`
	// Todo : overview chicken performance in (owner, kepala kandang)
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
	SellDate         string `json:"sellDate"`
	ChickenAge       uint64 `json:"chickenAge"`
	TotalSellChicken uint64 `json:"totalSellChicken"`
	Price            string `json:"price"`
	TotalPrice       string `json:"totalPrice"`
	PaymentStatus    string `json:"paymentStatus"`
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
	ChickenCage          ChickenCageResponse              `json:"chickenCage"`
	AfkirChickenCustomer AfkirChickenCustomerListResponse `json:"afkirChickenCustomer"`
	TotalSellChicken     uint64                           `json:"totalSellChicken"`
	PricePerChicken      string                           `json:"pricePerChicken"`
	TotalPrice           string                           `json:"totalPrice"`
}

type CreateAfkirChickenSaleRequest struct {
	AfkirChickenCustomerId  uint64                                `json:"afkirChickenCustomerId" validate:"required"`
	ChickenCageId           uint64                                `json:"chickenCageId" validate:"required"`
	TotalSellChicken        uint64                                `json:"totalSellChicken" validate:"required"`
	PricePerChicken         string                                `json:"pricePerChicken" validate:"required"`
	PaymentType             string                                `json:"paymentType" validate:"required,paymentType"`
	AfkirChickenSalePayment *CreateAfkirChickenSalePaymentRequest `json:"afkirChickenSalePayment"`
}

type AfkirChickenSaleResponse struct {
	Id                   uint64                           `json:"id"`
	AfkirChickenCustomer AfkirChickenCustomerListResponse `json:"afkirChickenCustomer"`
	ChickenCage          ChickenCageResponse              `json:"chickenCageId"`
	TotalSellChicken     uint64                           `json:"totalSellChicken"`
	PricePerChicken      string                           `json:"pricePerChicken"`
	TotalPrice           string                           `json:"totalPrice"`
	ChickenAge           string                           `json:"chickenAge"`
}

type CreateAfkirChickenSalePaymentRequest struct {
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
}

type UpdateAfkirChickenSalePaymentRequest struct {
	PaymentMethod string `json:"paymentMethod" validate:"required,paymentMethod"`
	PaymentDate   string `json:"paymentDate" validate:"required"`
	Nominal       string `json:"nominal" validate:"required,number"`
	PaymentProof  string `json:"paymentProof" validate:"required,url"`
}

type GetAfkirChickenSaleFilter struct {
	Page uint64 `query:"page"`
}
