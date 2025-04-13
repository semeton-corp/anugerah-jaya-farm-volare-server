package dto

type ChickenMonitoringRequest struct {
	CageId            uint64                            `json:"cageId" validate:"required"`
	ChickenCategory   string                            `json:"chickenCategory" validate:"required,chicken_category"`
	Age               uint64                            `json:"age" validate:"required,number"`
	TotalLiveChicken  uint64                            `json:"totalLiveChicken" validate:"required,number"`
	TotalSickChicken  uint64                            `json:"totalSickChicken" validate:"required,number"`
	TotalDeathChicken uint64                            `json:"totalDeathChicken" validate:"required,number"`
	TotalFeed         uint64                            `json:"totalFeed" validate:"required,number"`
	ChickenDisease    []ChickenDiseaseMonitoringRequest `json:"chickenDiseases" validate:"required"`
	ChickenVaccine    []ChickenVaccineMonitoringRequest `json:"chickenVaccines" validate:"required"`
}

type ChickenDiseaseMonitoringRequest struct {
	Disease  string  `json:"disease" validate:"required"`
	Medicine string  `json:"medicine" validate:"required"`
	Dose     float64 `json:"dose" validate:"required,number"`
	Unit     string  `json:"unit" validate:"required"`
}

type ChickenVaccineMonitoringRequest struct {
	Vaccine string  `json:"vaccine" validate:"required"`
	Dose    float64 `json:"dose" validate:"required"`
	Unit    string  `json:"unit" validate:"required"`
}

type ChickenMonitoringResponse struct {
	Id uint64 `json:"id"`
}
