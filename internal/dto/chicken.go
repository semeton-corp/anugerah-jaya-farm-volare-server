package dto

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type CreateChickenMonitoringRequest struct {
	ChickenCageId     uint64  `json:"cageId" validate:"required"`
	TotalSickChicken  uint64  `json:"totalSickChicken" validate:"required,number,min=0"`
	TotalDeathChicken uint64  `json:"totalDeathChicken" validate:"required,number,min=0"`
	TotalFeed         float64 `json:"totalFeed" validate:"required,number,min=0"`
	Note              string  `json:"note"`
}

type UpdateChickenMonitoringRequest struct {
	CageId            uint64                                  `json:"cageId" validate:"required"`
	ChickenCategory   string                                  `json:"chickenCategory" validate:"required,chicken_category"`
	Age               uint64                                  `json:"age" validate:"required,number"`
	TotalLiveChicken  uint64                                  `json:"totalLiveChicken" validate:"required,number,min=0"`
	TotalSickChicken  uint64                                  `json:"totalSickChicken" validate:"required,number,min=0"`
	TotalDeathChicken uint64                                  `json:"totalDeathChicken" validate:"required,number,min=0"`
	TotalFeed         float64                                 `json:"totalFeed" validate:"required,number,min=0"`
	ChickenDiseases   []UpdateChickenDiseaseMonitoringRequest `json:"chickenDiseases" validate:"required"`
	ChickenVaccines   []UpdateChickenVaccineMonitoringRequest `json:"chickenVaccines" validate:"required"`
}

type UpdateChickenDiseaseMonitoringRequest struct {
	Id       uint64  `json:"id"`
	Disease  string  `json:"disease" validate:"required"`
	Medicine string  `json:"medicine" validate:"required"`
	Dose     float64 `json:"dose" validate:"required,number"`
	Unit     string  `json:"unit" validate:"required"`
}

type UpdateChickenVaccineMonitoringRequest struct {
	Id      uint64  `json:"id"`
	Vaccine string  `json:"vaccine" validate:"required"`
	Dose    float64 `json:"dose" validate:"required"`
	Unit    string  `json:"unit" validate:"required"`
}

type CreateChickenDiseaseMonitoringRequest struct {
	Disease  string  `json:"disease" validate:"required"`
	Medicine string  `json:"medicine" validate:"required"`
	Dose     float64 `json:"dose" validate:"required,number"`
	Unit     string  `json:"unit" validate:"required"`
}

type CreateChickenVaccineMonitoringRequest struct {
	Vaccine string  `json:"vaccine" validate:"required"`
	Dose    float64 `json:"dose" validate:"required"`
	Unit    string  `json:"unit" validate:"required"`
}

type ChickenMonitoringResponse struct {
	Id                uint64                             `json:"id"`
	Cage              CageResponse                       `json:"cage"`
	Age               uint64                             `json:"age"`
	ChickenCategory   string                             `json:"chickenCategory"`
	TotalLiveChicken  uint64                             `json:"totalLiveChicken"`
	TotalSickChicken  uint64                             `json:"totalSickChicken"`
	TotalDeathChicken uint64                             `json:"totalDeathChicken"`
	TotalFeed         float64                            `json:"totalFeed"`
	ChickenDiseases   []ChickenDiseaseMonitoringResponse `json:"chickenDiseases"`
	ChickenVaccines   []ChickenVaccineMonitoringResponse `json:"chickenVaccines"`
}

type ChickenDiseaseMonitoringResponse struct {
	Id       uint64  `json:"id"`
	Disease  string  `json:"disease"`
	Medicine string  `json:"medicine"`
	Dose     float64 `json:"dose"`
	Unit     string  `json:"unit"`
}

type ChickenVaccineMonitoringResponse struct {
	Id      uint64  `json:"id"`
	Vaccine string  `json:"vaccine"`
	Dose    float64 `json:"dose"`
	Unit    string  `json:"unit"`
}

type ChickenMonitoringListResponse struct {
	Id                uint64                             `json:"id"`
	Cage              CageResponse                       `json:"cage"`
	ChickenCategory   string                             `json:"chickenCategory"`
	Age               uint64                             `json:"age"`
	TotalLiveChicken  uint64                             `json:"totalLiveChicken"`
	TotalSickChicken  uint64                             `json:"totalSickChicken"`
	TotalDeathChicken uint64                             `json:"totalDeathChicken"`
	TotalFeed         float64                            `json:"totalFeed"`
	MortalityRate     float64                            `json:"mortalityRate"`
	ChickenDiseases   []ChickenDiseaseMonitoringResponse `json:"chickenDiseases"`
	ChickenVaccines   []ChickenVaccineMonitoringResponse `json:"chickenVaccines"`
}

type GetChickenMonitoringFilter struct {
	Date      param.DateParam `query:"date"`
	StartDate param.DateParam
	EndDate   param.DateParam
	Location  uint64
}

type GetChickenOverviewFilter struct {
	Location          uint64                       `query:"location"`
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

type ChickenPieResponse struct {
	ChickenDOCType       float64 `json:"chickenDOCType"`
	ChickenGrowerType    float64 `json:"chickenGrowerType"`
	ChickentPreLayerType float64 `json:"chickentPreLayerType"`
	ChickenLayer         float64 `json:"chickenLayer"`
	ChickenAfkir         float64 `json:"chickenAfkir"`
}

type ChickenOverviewResponse struct {
	ChickenDetail ChickenDetailOverview  `json:"chickenDetail"`
	ChickenGraphs []ChickenGraphResponse `json:"chickenGraphs"`
	ChickenPie    ChickenPieResponse     `json:"chickenPie"`
}
