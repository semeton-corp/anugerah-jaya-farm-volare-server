package dto

type EggSummaryResponse struct {
	TotalGoodEggProductionInIkat   float64 `json:"totalGoodEggProductionInIkat"`
	TotalGoodEggProductionInKarpet float64 `json:"totalGoodEggProductionInKarpet"`
	TotalGoodEggProductionInButir  float64 `json:"totalGoodEggProductionInButir"`
	TotalGoodEggProductionInKg     float64 `json:"totalGoodEggProductionInKg"`
	TotalGoodEggSellInButir        float64 `json:"totalGoodEggSellInButir"`
	TotalGoodEggSellInKg           float64 `json:"totalGoodEggSellInKg"`
}

type ChickenSummaryResponse struct {
	TotalLiveChicken  uint64 `json:"totalLiveChicken"`
	TotalSickChicken  uint64 `json:"totalSickChicken"`
	TotalDeathChicken uint64 `json:"totalDeathChicken"`
}

type SaleSummaryResponse struct {
	Income string `json:"income"`
	Profit string `json:"profit"`
}

// Note : per day
type GeneralOverview struct {
	EggSummary           EggSummaryResponse           `json:"eggSummary"`
	WarehouseItemSummary WarehouseItemSummaryResponse `json:"warehouseItemSummary"`
	StoreItemSummary     StoreItemSummaryResponse     `json:"storeItemSummary"`
	SaleSummary          SaleSummaryResponse          `json:"saleSummary"`
	ChickenSummary       ChickenSummaryResponse       `json:"chickenSummary"`
}
