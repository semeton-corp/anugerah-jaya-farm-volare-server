package repository

import (
	"errors"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
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
	CountEggMonitoringByChickenCageIdToday(chickenCageId uint64) (int64, error)
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
	return r.GetDB().Model(&entity.EggMonitoring{}).Create(eggMonitoring).Error
}

func (r *EggRepository) GetEggMonitoringById(id uint64) (entity.EggMonitoring, error) {
	var eggMonitoring entity.EggMonitoring
	if err := r.GetDB().Model(entity.EggMonitoring{}).Preload("Warehouse.Location").Preload("ChickenCage.Cage.Location").Where("id = ?", id).First(&eggMonitoring).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.EggMonitoring{}, errx.NotFound("egg monitoring not found")
		}
		return entity.EggMonitoring{}, err
	}

	return eggMonitoring, nil
}

func (r *EggRepository) GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]entity.EggMonitoring, error) {
	eggMonitorings := make([]entity.EggMonitoring, 0)

	query := r.GetDB().
		Preload("Warehouse.Location").
		Preload("ChickenCage.Cage.Location").
		Model(&entity.EggMonitoring{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(egg_monitorings.created_at) = ?", filter.Date.Value())
	}

	if filter.LocationId > 0 {
		query = query.
			Joins("JOIN chicken_cages ON chicken_cages.id = egg_monitorings.chicken_cage_id").Joins("JOIN cages ON cages.id = chicken_cages.cage_id").
			Where("cages.location_id = ?", filter.LocationId)
	}

	if err := query.Order("egg_monitorings.created_at ASC").Find(&eggMonitorings).Error; err != nil {
		return nil, err
	}

	return eggMonitorings, nil
}

func (r *EggRepository) UpdateEggMonitoring(eggMonitoring *entity.EggMonitoring) error {
	return r.GetDB().Model(entity.EggMonitoring{}).Where("id = ?", eggMonitoring.Id).Updates(eggMonitoring).Error
}

func (r *EggRepository) DeleteEggMonitoring(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.EggMonitoring{}).Error
}

func (r *EggRepository) CountEggMonitoringByChickenCageIdToday(chickenCageId uint64) (int64, error) {
	var count int64
	if err := r.GetDB().Model(entity.EggMonitoring{}).Where("chicken_cage_id = ? AND DATE(created_at) = ?", chickenCageId, time.Now()).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
