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

func (r *CageRepository) CreateCage(data *entity.Cage) error {
	return r.GetDB().Create(data).Error
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
