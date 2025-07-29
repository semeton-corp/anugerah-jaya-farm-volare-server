package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type SaleOverviewSummaryResponse struct {
	Income            string `json:"income"`
	IsIncomeIncrease  bool   `json:"isIncomeIncrease"`
	IncomeNote        string `json:"incomeNote"`
	Profit            string `json:"profit"`
	IsProfitIncrease  bool   `json:"isProfitIncrease"`
	ProfitNote        string `json:"profitNote"`
	Expense           string `json:"expense"`
	IsExpenseIncrease bool   `json:"isExpenseIncrease"`
	ExpenseNote       string `json:"expenseNote"`
}

type EggSaleSummaryResponse struct {
	TotalGoodEggInKg        float64 `json:"totalGoodEggInKg"`
	TotalGoodEggInIkat      float64 `json:"totalGoodEggInIkat"`
	TotalCrackedEggInKg     float64 `json:"totalCrackedEggInKg"`
	TotalCrackedEggInIkat   float64 `json:"totalCrackedEggInIkat"`
	TotalBrokenEggInPlastik float64 `json:"totalBrokenEggInPlastik"`
}

type PlaceSaleSummaryResponse struct {
	Place      string `json:"place"`
	Income     string `json:"income"`
	Receivable string `json:"receivable"`
}

type SaleGraphResponse struct {
	Key     string `json:"key"`
	Income  string `json:"income"`
	Expense string `json:"expense"`
}

type EggSaleGraphResponse struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type LocationPieChartResponse struct {
	Place      string  `json:"place"`
	Percentage float64 `json:"percentage"` // Dari pendapatan
}

// Note : per month
type SaleOverviewResponse struct {
	SaleOverviewSummary SaleOverviewSummaryResponse `json:"saleOverviewSummary"`
	EggSaleSummary      EggSaleSummaryResponse      `json:"eggSaleSummary"`
	PlaceSaleSummary    []PlaceSaleSummaryResponse  `json:"placeSaleSummaries"`
	SaleGraphs          []SaleGraphResponse         `json:"saleGraphs"`
	EggSaleGraphs       []EggSaleGraphResponse      `json:"eggSaleGraphs"`
}

type GetSaleOverviewFilter struct {
	Month   param.MonthParam `query:"month"`
	Year    uint64           `query:"year"`
	EggType string           `query:"eggType" validate:"eggType"`
}
