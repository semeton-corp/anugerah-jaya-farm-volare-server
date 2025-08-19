package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type IncomeGraphResponse struct {
	Key        string  `json:"key"`
	Percentage float64 `json:"percentage"`
}

type IncomeListResponse struct {
	Date         string `json:"date"`
	PlaceName    string `json:"placeName"`
	Category     string `json:"category"`
	ItemName     string `json:"itemName"`
	ItemUnit     string `json:"itemUnit"`
	Quantity     string `json:"quantity"`
	CustomerName string `json:"customerName" `
	Nominal      string `json:"nominal"`
}

type IncomeOverviewResponse struct {
	IncomeGraphs []IncomeGraphResponse `json:"incomeGraphs"`
	Incomes      []IncomeListResponse  `json:"incomes"`
}

type GetIncomeOverviewFilter struct {
	Month      param.MonthParam `query:"month" validate:"required"`
	Year       uint64           `query:"year" validate:"required"`
	LocationId uint64           `query:"locationId"`
	Category   string           `query:"category"`
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
