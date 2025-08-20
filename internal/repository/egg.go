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
	GetEggMonitoringToday(chickenCageId uint64, date time.Time) (entity.EggMonitoring, error)
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
	if err := r.GetDB().Model(entity.EggMonitoring{}).Preload("Warehouse.Location").Preload("ChickenCage.Cage.Location").Preload("ChickenCage.Cage.CagePlacement.User.Role").Where("id = ?", id).First(&eggMonitoring).Error; err != nil {
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
		Preload("ChickenCage.Cage.CagePlacement.User.Role").
		Model(&entity.EggMonitoring{}).
		Joins("JOIN chicken_cages ON chicken_cages.id = egg_monitorings.chicken_cage_id").Joins("JOIN cages ON cages.id = chicken_cages.cage_id")

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(egg_monitorings.created_at) = ?", filter.Date.Value())
	}

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(egg_monitorings.created_at) >= ? AND DATE(egg_monitorings.created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.LocationId > 0 {
		query = query.Where("cages.location_id = ?", filter.LocationId)
	}

	if filter.CageId > 0 {
		query = query.Where("cages.id = ?", filter.CageId)
	}

	if err := query.Order("egg_monitorings.created_at ASC").Find(&eggMonitorings).Error; err != nil {
		return nil, err
	}

	return eggMonitorings, nil
}

func (r *EggRepository) UpdateEggMonitoring(eggMonitoring *entity.EggMonitoring) error {
	return r.GetDB().Model(&entity.EggMonitoring{}).Where("id = ?", eggMonitoring.Id).Updates(map[string]interface{}{
		"chicken_cage_id":          eggMonitoring.ChickenCageId,
		"warehouse_id":             eggMonitoring.WarehouseId,
		"total_weight_cracked_egg": eggMonitoring.TotalWeightCrackedEgg,
		"total_weight_good_egg":    eggMonitoring.TotalWeightGoodEgg,
		"total_good_egg":           eggMonitoring.TotalGoodEgg,
		"total_cracked_egg":        eggMonitoring.TotalCrackedEgg,
		"total_reject_egg":         eggMonitoring.TotalRejectEgg,
		"updated_by":               eggMonitoring.UpdatedBy,
	}).Error
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

func (r *EggRepository) GetEggMonitoringToday(chickenCageId uint64, date time.Time) (entity.EggMonitoring, error) {
	var monitoring entity.EggMonitoring

	err := r.GetDB().
		Where("chicken_cage_id = ? AND DATE(created_at) = ?", chickenCageId, date.Format("02 Jan 2006")).
		First(&monitoring).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.EggMonitoring{}, nil
		}
		return entity.EggMonitoring{}, err
	}

	return monitoring, nil
}
