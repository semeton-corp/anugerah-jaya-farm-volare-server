package dto

type EggSummaryResponse struct {
	TotalGoodEggProductionInIkat   float64 `json:"totalGoodEggProductionInIkat"`
	TotalGoodEggProductionInKarpet float64 `json:"totalGoodEggProductionInKarpet"`
	TotalGoodEggProductionInButir  float64 `json:"totalGoodEggProductionInButir"`
	TotalGoodEggProductionInKg     float64 `json:"totalGoodEggProductionInKg"`
	TotalGoodEggSaleInIkat         float64 `json:"totalGoodEggSaleInIkat"`
	TotalGoodEggSaleInKg           float64 `json:"totalGoodEggSaleInKg"`
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

type ProductionAndSaleEggGraphResponse struct {
	Key        string  `json:"key"`
	Production float64 `json:"production"`
	Sale       float64 `json:"sale"`
}

type GeneralOverview struct {
	EggSummary                 EggSummaryResponse                  `json:"eggSummary"`
	WarehouseItemSummary       WarehouseItemSummaryResponse        `json:"warehouseItemSummary"`
	StoreItemSummary           StoreItemSummaryResponse            `json:"storeItemSummary"`
	SaleSummary                SaleSummaryResponse                 `json:"saleSummary"`
	ChickenSummary             ChickenSummaryResponse              `json:"chickenSummary"`
	ProductionAndSaleEggGraphs []ProductionAndSaleEggGraphResponse `json:"productionAndSaleEggGraphs"`
}
