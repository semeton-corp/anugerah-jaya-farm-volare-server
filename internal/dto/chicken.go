package dto

import "time"

type CreateChickenMonitoringRequest struct {
	CageId            uint64                                  `json:"cageId" validate:"required"`
	ChickenCategory   string                                  `json:"chickenCategory" validate:"required,chicken_category"`
	Age               uint64                                  `json:"age" validate:"required,number"`
	TotalLiveChicken  uint64                                  `json:"totalLiveChicken" validate:"required,number"`
	TotalSickChicken  uint64                                  `json:"totalSickChicken" validate:"required,number"`
	TotalDeathChicken uint64                                  `json:"totalDeathChicken" validate:"required,number"`
	TotalFeed         float64                                 `json:"totalFeed" validate:"required,number"`
	ChickenDiseases   []CreateChickenDiseaseMonitoringRequest `json:"chickenDiseases" validate:"required"`
	ChickenVaccines   []CreateChickenVaccineMonitoringRequest `json:"chickenVaccines" validate:"required"`
}

type UpdateChickenMonitoringRequest struct {
	CageId            uint64  `json:"cageId" validate:"required"`
	ChickenCategory   string  `json:"chickenCategory" validate:"required,chicken_category"`
	Age               uint64  `json:"age" validate:"required,number"`
	TotalLiveChicken  uint64  `json:"totalLiveChicken" validate:"required,number"`
	TotalSickChicken  uint64  `json:"totalSickChicken" validate:"required,number"`
	TotalDeathChicken uint64  `json:"totalDeathChicken" validate:"required,number"`
	TotalFeed         float64 `json:"totalFeed" validate:"required,number"`
}

type UpdateChickenDiseaseMonitoringRequest struct {
	Disease  string  `json:"disease" validate:"required"`
	Medicine string  `json:"medicine" validate:"required"`
	Dose     float64 `json:"dose" validate:"required,number"`
	Unit     string  `json:"unit" validate:"required"`
}

type UpdateChickenVaccineMonitoringRequest struct {
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
	TotalLiveChicken  uint64                             `json:"totalLiveChicken"`
	TotalSickChicken  uint64                             `json:"totalSickChicken"`
	TotalDeathChicken uint64                             `json:"totalDeathChicken"`
	TotalFeed         float64                            `json:"totalFeed"`
	ChickenDiseases   []ChickenDiseaseMonitoringResponse `json:"chickenDiseases,omitempty"`
	ChickenVaccines   []ChickenVaccineMonitoringResponse `json:"chickenVaccines,omitempty"`
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
	Id                uint64       `json:"id"`
	Cage              CageResponse `json:"cage"`
	Age               uint64       `json:"age"`
	TotalLiveChicken  uint64       `json:"totalLiveChicken"`
	TotalSickChicken  uint64       `json:"totalSickChicken"`
	TotalDeathChicken uint64       `json:"totalDeathChicken"`
	TotalFeed         float64      `json:"totalFeed"`
	MortalityRate     float64      `json:"mortalityRate"`
}

type GetChickenMonitoringFilter struct {
	Date time.Time `query:"date"`
}
