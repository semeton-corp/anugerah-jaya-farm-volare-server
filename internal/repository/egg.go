package repository

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type EggRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IEggRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateEggMonitoring(eggMonitoring *entity.EggMonitoring) error
	GetEggMonitoringById(id uint64) (entity.EggMonitoring, error)
	GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]entity.EggMonitoring, error)
	UpdateEggMonitoring(eggMonitoring *entity.EggMonitoring) error
	DeleteEggMonitoring(id uint64) error
}

func NewEggRepository(db *gorm.DB) IEggRepository {
	return &EggRepository{
		db: db,
	}
}

func (r *EggRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *EggRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *EggRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *EggRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *EggRepository) CreateEggMonitoring(eggMonitoring *entity.EggMonitoring) error {
	return r.GetDB().Create(eggMonitoring).Error
}

func (r *EggRepository) GetEggMonitoringById(id uint64) (entity.EggMonitoring, error) {
	var eggMonitoring entity.EggMonitoring
	if err := r.GetDB().Model(entity.EggMonitoring{}).Preload("Cage.Location").Where("id = ?", id).First(&eggMonitoring).Error; err != nil {
		return entity.EggMonitoring{}, err
	}

	return eggMonitoring, nil
}

func (r *EggRepository) GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]entity.EggMonitoring, error) {
	var eggMonitorings []entity.EggMonitoring

	query := r.GetDB().Preload("Cage.Location").Model(&entity.EggMonitoring{})

	if !filter.Date.IsZero() {
		query = query.Where("created_at = ?", filter.Date)
	}

	if err := query.Find(&eggMonitorings).Order("created_at ASC").Error; err != nil {
		return nil, err
	}

	return eggMonitorings, nil
}

func (r *EggRepository) UpdateEggMonitoring(eggMonitoring *entity.EggMonitoring) error {
	return r.GetDB().Save(eggMonitoring).Error
}

func (r *EggRepository) DeleteEggMonitoring(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.EggMonitoring{}).Error
}
