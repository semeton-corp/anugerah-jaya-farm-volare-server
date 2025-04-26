package repository

import (
	"errors"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type ChickenRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IChickenRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateChickenMonitoring(chickenMonitoring *entity.ChickenMonitoring) error
	CreateChickenDiseaseMonitoring(chickenDisease *[]entity.ChickenDiseaseMonitoring) error
	CreateChickenVaccineMonitoring(chickenVaccine *[]entity.ChickenVaccineMonitoring) error
	GetChickenMonitoringById(id uint64) (entity.ChickenMonitoring, error)
	UpdateChickenMonitoring(chickenMonitoring *entity.ChickenMonitoring) error
	GetChickenMonitorings(filter *dto.GetChickenMonitoringFilter) ([]entity.ChickenMonitoring, error)
	GetChickenDiseaseMonitoringById(id uint64) (entity.ChickenDiseaseMonitoring, error)
	GetChickenVaccineMonitoringById(id uint64) (entity.ChickenVaccineMonitoring, error)
	UpdateChickenDiseaseMonitoring(chickenDiseaseMonitoring *entity.ChickenDiseaseMonitoring) error
	UpdateChickenVaccineMonitoring(chickenVaccineMonitoring *entity.ChickenVaccineMonitoring) error
	DeleteChickenMonitoring(id uint64) error
	DeleteChickenDiseaseMonitoring(id uint64) error
	DeleteChickenVaccineMonitoring(id uint64) error

	FirstOrCreateChickenVaccineMonitoring(chickenVaccineMonitoring *entity.ChickenVaccineMonitoring) error
	FirstOrCreateChickenDiseaseMonitoring(chickenDiseaseMonitoring *entity.ChickenDiseaseMonitoring) error

	SaveChickenVaccineMonitoring(chickenVaccineMonitoring *entity.ChickenVaccineMonitoring) error
	SaveChickenDiseaseMonitoring(chickenDiseaseMonitoring *entity.ChickenDiseaseMonitoring) error

	DeleteChickenVaccineMonitoringNotInIds(chickenMonitoringId uint64, keepIds []uint64) error
	DeleteChickenDiseaseMonitoringNotInIds(chickenMonitoringId uint64, keepIds []uint64) error

	CountChickenMonitoringByCageIdToday(cageId uint64) (int64, error)
}

func NewChickenRepository(db *gorm.DB) IChickenRepository {
	return &ChickenRepository{
		db: db,
	}
}

func (r *ChickenRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *ChickenRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *ChickenRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *ChickenRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *ChickenRepository) CountChickenMonitoringByCageIdToday(cageId uint64) (int64, error) {
	var count int64
	if err := r.GetDB().Model(entity.ChickenMonitoring{}).Where("cage_id = ? AND DATE(created_at) = ?", cageId, time.Now()).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ChickenRepository) CreateChickenMonitoring(chickenMonitoring *entity.ChickenMonitoring) error {
	return r.GetDB().Create(chickenMonitoring).Error
}

func (r *ChickenRepository) CreateChickenDiseaseMonitoring(chickenDisease *[]entity.ChickenDiseaseMonitoring) error {
	return r.GetDB().CreateInBatches(chickenDisease, len(*chickenDisease)).Error
}

func (r *ChickenRepository) CreateChickenVaccineMonitoring(chickenVaccine *[]entity.ChickenVaccineMonitoring) error {
	return r.GetDB().CreateInBatches(chickenVaccine, len(*chickenVaccine)).Error
}

func (r *ChickenRepository) GetChickenMonitoringById(id uint64) (entity.ChickenMonitoring, error) {
	var chickenMonitoring entity.ChickenMonitoring
	err := r.GetDB().
		Preload("Cage.Location").
		Preload("ChickenDiseaseMonitoring", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("ChickenVaccineMonitoring", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Where("id = ?", id).First(&chickenMonitoring).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ChickenMonitoring{}, errx.NotFound("chicken monitoring not found")
		}
		return entity.ChickenMonitoring{}, err
	}
	return chickenMonitoring, nil
}

func (r *ChickenRepository) UpdateChickenMonitoring(chickenMonitoring *entity.ChickenMonitoring) error {
	return r.GetDB().Model(entity.ChickenMonitoring{}).Where("id = ?", chickenMonitoring.Id).Updates(chickenMonitoring).Error
}

func (r *ChickenRepository) GetChickenMonitorings(filter *dto.GetChickenMonitoringFilter) ([]entity.ChickenMonitoring, error) {
	var chickenMonitorings []entity.ChickenMonitoring
	query := r.GetDB()

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Limit > 0 {
		query = query.Limit(int(filter.Limit))
	}

	err := query.
		Preload("ChickenDiseaseMonitoring").
		Preload("ChickenVaccineMonitoring").
		Preload("Cage.Location").
		Order("created_at desc").
		Find(&chickenMonitorings).Error

	if err != nil {
		return nil, err
	}
	return chickenMonitorings, nil
}

func (r *ChickenRepository) GetChickenDiseaseMonitoringById(id uint64) (entity.ChickenDiseaseMonitoring, error) {
	var chickenDiseaseMonitoring entity.ChickenDiseaseMonitoring
	err := r.GetDB().
		Where("id = ?", id).First(&chickenDiseaseMonitoring).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ChickenDiseaseMonitoring{}, errx.NotFound("chicken disease monitoring not found")
		}
		return entity.ChickenDiseaseMonitoring{}, err
	}
	return chickenDiseaseMonitoring, nil
}

func (r *ChickenRepository) GetChickenVaccineMonitoringById(id uint64) (entity.ChickenVaccineMonitoring, error) {
	var chickenVaccineMonitoring entity.ChickenVaccineMonitoring
	err := r.GetDB().
		Where("id = ?", id).First(&chickenVaccineMonitoring).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ChickenVaccineMonitoring{}, errx.NotFound("chicken vaccine monitoring not found")
		}
		return entity.ChickenVaccineMonitoring{}, err
	}
	return chickenVaccineMonitoring, nil
}

func (r *ChickenRepository) UpdateChickenDiseaseMonitoring(chickenDiseaseMonitoring *entity.ChickenDiseaseMonitoring) error {
	return r.GetDB().Model(entity.ChickenDiseaseMonitoring{}).Where("id = ?", chickenDiseaseMonitoring.Id).Updates(chickenDiseaseMonitoring).Error
}

func (r *ChickenRepository) UpdateChickenVaccineMonitoring(chickenVaccineMonitoring *entity.ChickenVaccineMonitoring) error {
	return r.GetDB().Model(entity.ChickenVaccineMonitoring{}).Where("id = ?", chickenVaccineMonitoring.Id).Updates(chickenVaccineMonitoring).Error
}

func (r *ChickenRepository) DeleteChickenMonitoring(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ChickenMonitoring{}).Error
}

func (r *ChickenRepository) DeleteChickenDiseaseMonitoring(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ChickenDiseaseMonitoring{}).Error
}

func (r *ChickenRepository) DeleteChickenVaccineMonitoring(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ChickenVaccineMonitoring{}).Error
}

func (r *ChickenRepository) FirstOrCreateChickenVaccineMonitoring(chickenVaccineMonitoring *entity.ChickenVaccineMonitoring) error {
	return r.GetDB().FirstOrCreate(chickenVaccineMonitoring, entity.ChickenVaccineMonitoring{
		Id: chickenVaccineMonitoring.Id,
	}).Error
}

func (r *ChickenRepository) FirstOrCreateChickenDiseaseMonitoring(chickenDiseaseMonitoring *entity.ChickenDiseaseMonitoring) error {
	return r.GetDB().FirstOrCreate(chickenDiseaseMonitoring, entity.ChickenDiseaseMonitoring{
		Id: chickenDiseaseMonitoring.Id,
	}).Error
}

func (r *ChickenRepository) SaveChickenVaccineMonitoring(chickenVaccineMonitoring *entity.ChickenVaccineMonitoring) error {
	return r.GetDB().Save(chickenVaccineMonitoring).Error
}

func (r *ChickenRepository) SaveChickenDiseaseMonitoring(chickenDiseaseMonitoring *entity.ChickenDiseaseMonitoring) error {
	return r.GetDB().Save(chickenDiseaseMonitoring).Error
}

func (r *ChickenRepository) DeleteChickenVaccineMonitoringNotInIds(chickenMonitoringId uint64, keepIds []uint64) error {
	return r.GetDB().
		Where("chicken_monitoring_id = ?", chickenMonitoringId).
		Where("id NOT IN ?", keepIds).
		Delete(&entity.ChickenVaccineMonitoring{}).Error
}

func (r *ChickenRepository) DeleteChickenDiseaseMonitoringNotInIds(chickenMonitoringId uint64, keepIds []uint64) error {
	return r.GetDB().
		Where("chicken_monitoring_id = ?", chickenMonitoringId).
		Where("id NOT IN ?", keepIds).
		Delete(&entity.ChickenDiseaseMonitoring{}).Error
}
