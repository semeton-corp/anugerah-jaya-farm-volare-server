package repository

import (
	"errors"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type CageRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type ICageRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetCages(filter dto.GetCageFilter) ([]entity.Cage, error)
	CreateCage(data *entity.Cage) error
	GetCageById(id uint64) (entity.Cage, error)
	UpdateCage(cage *entity.Cage) error
	DeleteCage(id uint64) error
	GetCageByIds(ids []uint64) ([]entity.Cage, error)

	CreateChickenCage(chickenCage *entity.ChickenCage) error
	GetChickenCages(filter dto.GetChickenCageFilter) ([]entity.ChickenCage, error)
	GetChickenCageById(id uint64) (entity.ChickenCage, error)
}

func NewCageRepository(db *gorm.DB) ICageRepository {
	return &CageRepository{
		db: db,
	}
}

func (r *CageRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *CageRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *CageRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *CageRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *CageRepository) GetCages(filter dto.GetCageFilter) ([]entity.Cage, error) {
	var cages []entity.Cage
	query := r.GetDB()

	if filter.LocationId > 0 {
		query.Where("location_id = ?", filter.LocationId)
	}

	err := query.Preload("Location").Find(&cages).Error
	if err != nil {
		return nil, err
	}

	return cages, nil
}

func (r *CageRepository) CreateCage(cage *entity.Cage) error {
	return r.GetDB().Create(&cage).Error
}

func (r *CageRepository) GetCageById(id uint64) (entity.Cage, error) {
	var cage entity.Cage
	if err := r.GetDB().Preload("Location").Where("id = ?", id).First(&cage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Cage{}, errx.NotFound("cage not found")
		}
		return entity.Cage{}, err
	}

	return cage, nil
}

func (r *CageRepository) UpdateCage(cage *entity.Cage) error {
	return r.GetDB().Model(&entity.Cage{}).Where("id = ?", cage.Id).Updates(&cage).Error
}

func (r *CageRepository) DeleteCage(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Cage{}).Error
}

func (r *CageRepository) GetChickenCageByCageId(cageId uint64) (entity.ChickenCage, error) {
	var chickenCage entity.ChickenCage
	// Find the newest chicken in the cage
	err := r.GetDB().Model(&entity.ChickenCage{}).Where("cage_id = ?", cageId).Order("created_at  DESC").First(&chickenCage).Error

	if err != nil {
		return entity.ChickenCage{}, err
	}

	return chickenCage, nil
}

func (r *CageRepository) CreateChickenCage(chickenCage *entity.ChickenCage) error {
	return r.GetDB().Model(&entity.ChickenCage{}).Create(&chickenCage).Error
}

func (r *CageRepository) GetChickenCages(filter dto.GetChickenCageFilter) ([]entity.ChickenCage, error) {
	var chickenCages []entity.ChickenCage
	query := r.GetDB().Model(&entity.ChickenCage{})

	if filter.LocationId > 0 {
		query = query.Joins("JOIN cages ON cages.id = chicken_cages.cage_id").
			Where("cages.location_id = ?", filter.LocationId)
	}

	// Note : Subquery to get the newest chicken_cage per cage_id
	subQuery := r.GetDB().Model(&entity.ChickenCage{}).
		Select("MAX(id)").
		Group("cage_id")

	query = query.Where("chicken_cages.id IN (?)", subQuery)

	err := query.
		Preload("Cage.Location").
		Preload("ChickenProcurement").
		Preload("Cage.CagePlacement.User.Role").
		Order("chicken_cages.created_at DESC").
		Find(&chickenCages).Error

	if err != nil {
		return nil, err
	}

	return chickenCages, nil
}

func (r *CageRepository) GetChickenCageById(id uint64) (entity.ChickenCage, error) {
	var chickenCage entity.ChickenCage
	err := r.GetDB().Preload("Cage.Location").Preload("ChickenProcurement").Preload("Cage.CagePlacement.User.Role").Where("id = ?", id).First(&chickenCage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ChickenCage{}, errx.NotFound("chicken cage not found")
		}
		return entity.ChickenCage{}, err
	}

	return chickenCage, nil
}

func (r *CageRepository) GetCageByIds(ids []uint64) ([]entity.Cage, error) {
	var cages []entity.Cage
	err := r.GetDB().Model(&entity.Cage{}).Where("id IN ?", ids).Find(&cages).Error
	if err != nil {
		return nil, err
	}

	return cages, nil
}
