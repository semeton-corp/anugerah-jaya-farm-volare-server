package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

type ChickenService struct {
	log        *zap.Logger
	repository repository.IChickenRepository
}

type IChickenService interface {
	CreateChickenMonitoring(request dto.CreateChickenMonitoringRequest, accoundId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	GetChickenMonitorings(filter dto.GetChickenMonitoringFilter) ([]dto.ChickenMonitoringListResponse, error)
	GetChickenMonitoringById(id uint64) (dto.ChickenMonitoringResponse, error)
	UpdateChickenMonitoring(id uint64, request dto.UpdateChickenMonitoringRequest, accoundId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	DeleteChickenMonitoring(id uint64) error

	CreateChickenDiseaseMonitoring(chickenMonitoringId uint64, request dto.CreateChickenDiseaseMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	UpdateChickenDiseaseMonitoring(id uint64, request dto.UpdateChickenDiseaseMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	DeleteChickenDiseaseMonitoring(id uint64) error

	CreateChickenVaccineMonitoring(chickenMonitoringId uint64, request dto.CreateChickenVaccineMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	UpdateChickenVaccineMonitoring(id uint64, request dto.UpdateChickenVaccineMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	DeleteChickenVaccineMonitoring(id uint64) error
}

func NewChickenService(log *zap.Logger, repository repository.IChickenRepository) IChickenService {
	return &ChickenService{
		log:        log,
		repository: repository,
	}
}

func (c *ChickenService) CreateChickenMonitoring(request dto.CreateChickenMonitoringRequest, accoundId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	c.repository.UseTx(true)
	defer c.repository.Rollback()

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		return dto.ChickenMonitoringResponse{}, errx.BadRequest("invalid chicken category")
	}

	chickenMonitoring := entity.ChickenMonitoring{
		CageId:            request.CageId,
		Age:               request.Age,
		ChickenCategory:   chickenCategory,
		TotalLiveChicken:  request.TotalLiveChicken,
		TotalDeathChicken: request.TotalDeathChicken,
		TotalSickChicken:  request.TotalSickChicken,
		TotalFeed:         request.TotalFeed,
		CreatedBy:         accoundId,
	}

	err := c.repository.CreateChickenMonitoring(&chickenMonitoring)
	if err != nil {
		c.log.Error("[CreateChickenMonitoring] failed to create chicken monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	if request.ChickenDiseases != nil {
		chickenDiseases := make([]entity.ChickenDiseaseMonitoring, len(request.ChickenDiseases))
		for i, disease := range request.ChickenDiseases {
			chickenDiseases[i] = entity.ChickenDiseaseMonitoring{
				ChickenMonitoringId: chickenMonitoring.Id,
				Disease:             disease.Disease,
				Medicine:            disease.Medicine,
				Dose:                disease.Dose,
				Unit:                disease.Unit,
				CreatedBy:           accoundId,
			}
		}

		err = c.repository.CreateChickenDiseaseMonitoring(&chickenDiseases)
		if err != nil {
			c.log.Error("[CreateChickenMonitoring] failed to create chicken diseases", zap.Error(err))
			return dto.ChickenMonitoringResponse{}, err
		}
	}

	if request.ChickenVaccines != nil {
		chickenVaccine := make([]entity.ChickenVaccineMonitoring, len(request.ChickenVaccines))
		for i, vaccine := range request.ChickenVaccines {
			chickenVaccine[i] = entity.ChickenVaccineMonitoring{
				ChickenMonitoringId: chickenMonitoring.Id,
				Vaccine:             vaccine.Vaccine,
				Dose:                vaccine.Dose,
				Unit:                vaccine.Unit,
				CreatedBy:           accoundId,
			}
		}

		err = c.repository.CreateChickenVaccineMonitoring(&chickenVaccine)
		if err != nil {
			c.log.Error("[CreateChickenMonitoring] failed to create chicken vaccines", zap.Error(err))
			return dto.ChickenMonitoringResponse{}, err
		}
	}

	err = c.repository.Commit()
	if err != nil {
		c.log.Error("[CreateChickenMonitoring] failed to commit transaction", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err = c.repository.GetChickenMonitoringById(chickenMonitoring.Id)
	if err != nil {
		c.log.Error("[CreateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenDiseasesResponse := make([]dto.ChickenDiseaseMonitoringResponse, len(chickenMonitoring.ChickenDiseaseMonitoring))
	for i, disease := range chickenMonitoring.ChickenDiseaseMonitoring {
		chickenDiseasesResponse[i] = dto.ChickenDiseaseMonitoringResponse{
			Id:       disease.Id,
			Disease:  disease.Disease,
			Medicine: disease.Medicine,
			Dose:     disease.Dose,
			Unit:     disease.Unit,
		}
	}

	chickenVaccinesResponse := make([]dto.ChickenVaccineMonitoringResponse, len(chickenMonitoring.ChickenVaccineMonitoring))
	for i, vaccine := range chickenMonitoring.ChickenVaccineMonitoring {
		chickenVaccinesResponse[i] = dto.ChickenVaccineMonitoringResponse{
			Id:      vaccine.Id,
			Vaccine: vaccine.Vaccine,
			Dose:    vaccine.Dose,
			Unit:    vaccine.Unit,
		}
	}

	return dto.ChickenMonitoringResponse{
		Id:              chickenMonitoring.Id,
		ChickenCategory: chickenCategory.String(),
		Cage: dto.CageResponse{
			Id:   chickenMonitoring.Cage.Id,
			Name: chickenMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				Id:   chickenMonitoring.Cage.Location.Id,
				Name: chickenMonitoring.Cage.Location.Name,
			},
		},
		Age:               chickenMonitoring.Age,
		TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		ChickenDiseases:   chickenDiseasesResponse,
		ChickenVaccines:   chickenVaccinesResponse,
	}, nil
}

func (c *ChickenService) GetChickenMonitoringById(id uint64) (dto.ChickenMonitoringResponse, error) {
	chickenMonitoring, err := c.repository.GetChickenMonitoringById(id)
	if err != nil {
		c.log.Error("[GetChickenMonitoringById] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenDiseasesResponse := make([]dto.ChickenDiseaseMonitoringResponse, len(chickenMonitoring.ChickenDiseaseMonitoring))
	for i, disease := range chickenMonitoring.ChickenDiseaseMonitoring {
		chickenDiseasesResponse[i] = dto.ChickenDiseaseMonitoringResponse{
			Id:       disease.Id,
			Disease:  disease.Disease,
			Medicine: disease.Medicine,
			Dose:     disease.Dose,
			Unit:     disease.Unit,
		}
	}

	chickenVaccinesResponse := make([]dto.ChickenVaccineMonitoringResponse, len(chickenMonitoring.ChickenVaccineMonitoring))
	for i, vaccine := range chickenMonitoring.ChickenVaccineMonitoring {
		chickenVaccinesResponse[i] = dto.ChickenVaccineMonitoringResponse{
			Id:      vaccine.Id,
			Vaccine: vaccine.Vaccine,
			Dose:    vaccine.Dose,
			Unit:    vaccine.Unit,
		}
	}

	return dto.ChickenMonitoringResponse{
		Id: chickenMonitoring.Id,
		Cage: dto.CageResponse{
			Id:   chickenMonitoring.Cage.Id,
			Name: chickenMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				Id:   chickenMonitoring.Cage.Location.Id,
				Name: chickenMonitoring.Cage.Location.Name,
			},
		},
		Age:               chickenMonitoring.Age,
		TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		ChickenDiseases:   chickenDiseasesResponse,
		ChickenVaccines:   chickenVaccinesResponse,
	}, nil
}

func (c *ChickenService) GetChickenMonitorings(filter dto.GetChickenMonitoringFilter) ([]dto.ChickenMonitoringListResponse, error) {
	chickenMonitorings, err := c.repository.GetChickenMonitorings(&filter)
	if err != nil {
		c.log.Error("[GetChickenMonitorings] failed to get chicken monitorings", zap.Error(err))
		return []dto.ChickenMonitoringListResponse{}, err
	}

	chickenMonitoringsResponse := make([]dto.ChickenMonitoringListResponse, len(chickenMonitorings))
	for i, chickenMonitoring := range chickenMonitorings {
		chickenMonitoringsResponse[i] = dto.ChickenMonitoringListResponse{
			Id:              chickenMonitoring.Id,
			ChickenCategory: chickenMonitoring.ChickenCategory.String(),
			Cage: dto.CageResponse{
				Id:   chickenMonitoring.Cage.Id,
				Name: chickenMonitoring.Cage.Name,
				Location: dto.LocationResponse{
					Id:   chickenMonitoring.Cage.Location.Id,
					Name: chickenMonitoring.Cage.Location.Name,
				},
			},
			Age:               chickenMonitoring.Age,
			TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
			TotalSickChicken:  chickenMonitoring.TotalSickChicken,
			TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
			TotalFeed:         chickenMonitoring.TotalFeed,
			MortalityRate:     float64((chickenMonitoring.TotalDeathChicken / (chickenMonitoring.TotalLiveChicken + chickenMonitoring.TotalSickChicken)) * 100.0),
		}
	}

	return chickenMonitoringsResponse, nil
}

func (c *ChickenService) UpdateChickenMonitoring(id uint64, request dto.UpdateChickenMonitoringRequest, accoundId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	c.repository.UseTx(false)
	chickenMonitoring, err := c.repository.GetChickenMonitoringById(id)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring.CageId = request.CageId
	chickenMonitoring.Age = request.Age
	chickenMonitoring.TotalLiveChicken = request.TotalLiveChicken
	chickenMonitoring.TotalSickChicken = request.TotalSickChicken
	chickenMonitoring.TotalDeathChicken = request.TotalDeathChicken
	chickenMonitoring.TotalFeed = request.TotalFeed
	chickenMonitoring.UpdateBy = accoundId

	err = c.repository.UpdateChickenMonitoring(&chickenMonitoring)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to update chicken monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err = c.repository.GetChickenMonitoringById(chickenMonitoring.Id)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenDiseasesResponse := make([]dto.ChickenDiseaseMonitoringResponse, len(chickenMonitoring.ChickenDiseaseMonitoring))
	for i, disease := range chickenMonitoring.ChickenDiseaseMonitoring {
		chickenDiseasesResponse[i] = dto.ChickenDiseaseMonitoringResponse{
			Id:       disease.Id,
			Disease:  disease.Disease,
			Medicine: disease.Medicine,
			Dose:     disease.Dose,
			Unit:     disease.Unit,
		}
	}

	chickenVaccinesResponse := make([]dto.ChickenVaccineMonitoringResponse, len(chickenMonitoring.ChickenVaccineMonitoring))
	for i, vaccine := range chickenMonitoring.ChickenVaccineMonitoring {
		chickenVaccinesResponse[i] = dto.ChickenVaccineMonitoringResponse{
			Id:      vaccine.Id,
			Vaccine: vaccine.Vaccine,
			Dose:    vaccine.Dose,
			Unit:    vaccine.Unit,
		}
	}

	return dto.ChickenMonitoringResponse{
		Id:              chickenMonitoring.Id,
		ChickenCategory: chickenMonitoring.ChickenCategory.String(),
		Cage: dto.CageResponse{
			Id:   chickenMonitoring.Cage.Id,
			Name: chickenMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				Id:   chickenMonitoring.Cage.Location.Id,
				Name: chickenMonitoring.Cage.Location.Name,
			},
		},
		Age:               chickenMonitoring.Age,
		TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		ChickenDiseases:   chickenDiseasesResponse,
		ChickenVaccines:   chickenVaccinesResponse,
	}, nil
}

func (c *ChickenService) CreateChickenDiseaseMonitoring(chickenMonitoringId uint64, request dto.CreateChickenDiseaseMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	chickenDisease := entity.ChickenDiseaseMonitoring{
		ChickenMonitoringId: chickenMonitoringId,
		Disease:             request.Disease,
		Medicine:            request.Medicine,
		Dose:                request.Dose,
		Unit:                request.Unit,
		CreatedBy:           accountId,
	}

	err := c.repository.CreateChickenDiseaseMonitoring(&[]entity.ChickenDiseaseMonitoring{chickenDisease})
	if err != nil {
		c.log.Error("[CreateChickenDiseaseMonitoring] failed to create chicken disease monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err := c.repository.GetChickenMonitoringById(chickenMonitoringId)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenDiseasesResponse := make([]dto.ChickenDiseaseMonitoringResponse, len(chickenMonitoring.ChickenDiseaseMonitoring))
	for i, disease := range chickenMonitoring.ChickenDiseaseMonitoring {
		chickenDiseasesResponse[i] = dto.ChickenDiseaseMonitoringResponse{
			Id:       disease.Id,
			Disease:  disease.Disease,
			Medicine: disease.Medicine,
			Dose:     disease.Dose,
			Unit:     disease.Unit,
		}
	}

	chickenVaccinesResponse := make([]dto.ChickenVaccineMonitoringResponse, len(chickenMonitoring.ChickenVaccineMonitoring))
	for i, vaccine := range chickenMonitoring.ChickenVaccineMonitoring {
		chickenVaccinesResponse[i] = dto.ChickenVaccineMonitoringResponse{
			Id:      vaccine.Id,
			Vaccine: vaccine.Vaccine,
			Dose:    vaccine.Dose,
			Unit:    vaccine.Unit,
		}
	}

	return dto.ChickenMonitoringResponse{
		Id:              chickenMonitoring.Id,
		ChickenCategory: chickenMonitoring.ChickenCategory.String(),
		Cage: dto.CageResponse{
			Id:   chickenMonitoring.Cage.Id,
			Name: chickenMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				Id:   chickenMonitoring.Cage.Location.Id,
				Name: chickenMonitoring.Cage.Location.Name,
			},
		},
		Age:               chickenMonitoring.Age,
		TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		ChickenDiseases:   chickenDiseasesResponse,
		ChickenVaccines:   chickenVaccinesResponse,
	}, nil
}

func (c *ChickenService) CreateChickenVaccineMonitoring(chickenMonitoringId uint64, request dto.CreateChickenVaccineMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	chickenVaccine := entity.ChickenVaccineMonitoring{
		ChickenMonitoringId: chickenMonitoringId,
		Vaccine:             request.Vaccine,
		Dose:                request.Dose,
		Unit:                request.Unit,
		CreatedBy:           accountId,
	}

	err := c.repository.CreateChickenVaccineMonitoring(&[]entity.ChickenVaccineMonitoring{chickenVaccine})
	if err != nil {
		c.log.Error("[CreateChickenVaccineMonitoring] failed to create chicken vaccine monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err := c.repository.GetChickenMonitoringById(chickenMonitoringId)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenDiseasesResponse := make([]dto.ChickenDiseaseMonitoringResponse, len(chickenMonitoring.ChickenDiseaseMonitoring))
	for i, disease := range chickenMonitoring.ChickenDiseaseMonitoring {
		chickenDiseasesResponse[i] = dto.ChickenDiseaseMonitoringResponse{
			Id:       disease.Id,
			Disease:  disease.Disease,
			Medicine: disease.Medicine,
			Dose:     disease.Dose,
			Unit:     disease.Unit,
		}
	}

	chickenVaccinesResponse := make([]dto.ChickenVaccineMonitoringResponse, len(chickenMonitoring.ChickenVaccineMonitoring))
	for i, vaccine := range chickenMonitoring.ChickenVaccineMonitoring {
		chickenVaccinesResponse[i] = dto.ChickenVaccineMonitoringResponse{
			Id:      vaccine.Id,
			Vaccine: vaccine.Vaccine,
			Dose:    vaccine.Dose,
			Unit:    vaccine.Unit,
		}
	}

	return dto.ChickenMonitoringResponse{
		Id:              chickenMonitoring.Id,
		ChickenCategory: chickenMonitoring.ChickenCategory.String(),
		Cage: dto.CageResponse{
			Id:   chickenMonitoring.Cage.Id,
			Name: chickenMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				Id:   chickenMonitoring.Cage.Location.Id,
				Name: chickenMonitoring.Cage.Location.Name,
			},
		},
		Age:               chickenMonitoring.Age,
		TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		ChickenDiseases:   chickenDiseasesResponse,
		ChickenVaccines:   chickenVaccinesResponse,
	}, nil
}

func (c *ChickenService) UpdateChickenDiseaseMonitoring(id uint64, request dto.UpdateChickenDiseaseMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	chickenDisease, err := c.repository.GetChickenDiseaseMonitoringById(id)
	if err != nil {
		c.log.Error("[UpdateChickenDiseaseMonitoring] failed to get chicken disease monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenDisease.Disease = request.Disease
	chickenDisease.Medicine = request.Medicine
	chickenDisease.Dose = request.Dose
	chickenDisease.Unit = request.Unit
	chickenDisease.UpdatedBy = accountId

	err = c.repository.UpdateChickenDiseaseMonitoring(&chickenDisease)
	if err != nil {
		c.log.Error("[UpdateChickenDiseaseMonitoring] failed to update chicken disease monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err := c.repository.GetChickenMonitoringById(chickenDisease.ChickenMonitoringId)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenDiseasesResponse := make([]dto.ChickenDiseaseMonitoringResponse, len(chickenMonitoring.ChickenDiseaseMonitoring))
	for i, disease := range chickenMonitoring.ChickenDiseaseMonitoring {
		chickenDiseasesResponse[i] = dto.ChickenDiseaseMonitoringResponse{
			Id:       disease.Id,
			Disease:  disease.Disease,
			Medicine: disease.Medicine,
			Dose:     disease.Dose,
			Unit:     disease.Unit,
		}
	}

	chickenVaccinesResponse := make([]dto.ChickenVaccineMonitoringResponse, len(chickenMonitoring.ChickenVaccineMonitoring))
	for i, vaccine := range chickenMonitoring.ChickenVaccineMonitoring {
		chickenVaccinesResponse[i] = dto.ChickenVaccineMonitoringResponse{
			Id:      vaccine.Id,
			Vaccine: vaccine.Vaccine,
			Dose:    vaccine.Dose,
			Unit:    vaccine.Unit,
		}
	}

	return dto.ChickenMonitoringResponse{
		Id:              chickenMonitoring.Id,
		ChickenCategory: chickenMonitoring.ChickenCategory.String(),
		Cage: dto.CageResponse{
			Id:   chickenMonitoring.Cage.Id,
			Name: chickenMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				Id:   chickenMonitoring.Cage.Location.Id,
				Name: chickenMonitoring.Cage.Location.Name,
			},
		},
		Age:               chickenMonitoring.Age,
		TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		ChickenDiseases:   chickenDiseasesResponse,
		ChickenVaccines:   chickenVaccinesResponse,
	}, nil
}

func (c *ChickenService) UpdateChickenVaccineMonitoring(id uint64, request dto.UpdateChickenVaccineMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	chickenVaccine, err := c.repository.GetChickenVaccineMonitoringById(id)
	if err != nil {
		c.log.Error("[UpdateChickenVaccineMonitoring] failed to get chicken vaccine monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenVaccine.Vaccine = request.Vaccine
	chickenVaccine.Dose = request.Dose
	chickenVaccine.Unit = request.Unit
	chickenVaccine.UpdatedBy = accountId

	err = c.repository.UpdateChickenVaccineMonitoring(&chickenVaccine)
	if err != nil {
		c.log.Error("[UpdateChickenVaccineMonitoring] failed to update chicken vaccine monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err := c.repository.GetChickenMonitoringById(chickenVaccine.ChickenMonitoringId)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenDiseasesResponse := make([]dto.ChickenDiseaseMonitoringResponse, len(chickenMonitoring.ChickenDiseaseMonitoring))
	for i, disease := range chickenMonitoring.ChickenDiseaseMonitoring {
		chickenDiseasesResponse[i] = dto.ChickenDiseaseMonitoringResponse{
			Id:       disease.Id,
			Disease:  disease.Disease,
			Medicine: disease.Medicine,
			Dose:     disease.Dose,
			Unit:     disease.Unit,
		}
	}

	chickenVaccinesResponse := make([]dto.ChickenVaccineMonitoringResponse, len(chickenMonitoring.ChickenVaccineMonitoring))
	for i, vaccine := range chickenMonitoring.ChickenVaccineMonitoring {
		chickenVaccinesResponse[i] = dto.ChickenVaccineMonitoringResponse{
			Id:      vaccine.Id,
			Vaccine: vaccine.Vaccine,
			Dose:    vaccine.Dose,
			Unit:    vaccine.Unit,
		}
	}

	return dto.ChickenMonitoringResponse{
		Id:              chickenMonitoring.Id,
		ChickenCategory: chickenMonitoring.ChickenCategory.String(),
		Cage: dto.CageResponse{
			Id:   chickenMonitoring.Cage.Id,
			Name: chickenMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				Id:   chickenMonitoring.Cage.Location.Id,
				Name: chickenMonitoring.Cage.Location.Name,
			},
		},
		Age:               chickenMonitoring.Age,
		TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		ChickenDiseases:   chickenDiseasesResponse,
		ChickenVaccines:   chickenVaccinesResponse,
	}, nil
}

func (c *ChickenService) DeleteChickenMonitoring(id uint64) error {
	err := c.repository.DeleteChickenMonitoring(id)
	if err != nil {
		c.log.Error("[DeleteChickenMonitoring] failed to delete chicken monitoring", zap.Error(err))
		return err
	}

	return nil
}

func (c *ChickenService) DeleteChickenDiseaseMonitoring(id uint64) error {
	err := c.repository.DeleteChickenDiseaseMonitoring(id)
	if err != nil {
		c.log.Error("[DeleteChickenDiseaseMonitoring] failed to delete chicken disease monitoring", zap.Error(err))
		return err
	}

	return nil
}

func (c *ChickenService) DeleteChickenVaccineMonitoring(id uint64) error {
	err := c.repository.DeleteChickenVaccineMonitoring(id)
	if err != nil {
		c.log.Error("[DeleteChickenVaccineMonitoring] failed to delete chicken vaccine monitoring", zap.Error(err))
		return err
	}

	return nil
}
