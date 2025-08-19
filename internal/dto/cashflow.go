package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type IncomePieResponse struct {
	WarehouseEggSalePercentage float64 `json:"warehouseEggSalePercentage"`
	StoreEggSalePercentage     float64 `json:"storeEggSalePercentage"`
	AfkirChickenSalePercentage float64 `json:"afkirChickenSalePercentage"`
}

type IncomeListResponse struct {
	ParentId     uint64 `json:"parentId"`
	Id           uint64 `json:"id"`
	Date         string `json:"date"`
	PlaceName    string `json:"placeName"`
	Category     string `json:"category"`
	ItemName     string `json:"itemName"`
	ItemUnit     string `json:"itemUnit"`
	Quantity     string `json:"quantity"`
	CustomerName string `json:"customerName" `
	Nominal      string `json:"nominal"`
	PaymentProof string `json:"paymentProof"`
}

type IncomeResponse struct {
	Id                  uint64 `json:"id"`
	Date                string `json:"date"`
	Time                string `json:"time"`
	Category            string `json:"category"`
	PlaceName           string `json:"placeName"`
	CustomerName        string `json:"customerName" `
	CustomerPhoneNumber string `json:"customerPhoneNumber"`
	ItemName            string `json:"itemName"`
	ItemUnit            string `json:"itemUnit"`
	Quantity            string `json:"quantity"`
	Nominal             string `json:"nominal"`
	PaymentType         string `json:"paymentType"`
	TotalPrice          string `json:"totalPrice"`
	PaymentMethod       string `json:"paymentMethod"`
	InputBy             string `json:"inputBy"`
	PaymentProof        string `json:"paymentProof"`
}

type IncomeOverviewResponse struct {
	IncomePies IncomePieResponse    `json:"incomePies"`
	Incomes    []IncomeListResponse `json:"incomes"`
}

type GetIncomeOverviewFilter struct {
	Month          param.MonthParam `query:"month" validate:"required"`
	Year           uint64           `query:"year" validate:"required"`
	IncomeCategory string           `query:"category" validate:"required,incomeCategory"`
}

type ExpenseGraphResponse struct {
	Key        string  `json:"key"`
	Percentage float64 `json:"percentage"`
}

type ExpenseListResponse struct {
}

type ExpenseOverviewResponse struct {
	ExpenseGraphs []ExpenseGraphResponse `json:"expenseGraphs"`
}

type GetExpenseOverviewFilter struct {
}

type GetReceivablesOverviewFilter struct {
}

type GetDebtOverviewFilter struct {
}

type GetSaleCashflowFilter struct {
	Year  uint64           `query:"year"`
	Month param.MonthParam `query:"month"`
}
