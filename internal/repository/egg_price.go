package repository

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type EggPriceRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IEggPriceRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateEggPrice(eggPrice *entity.EggPrice) error
	GetEggPrices() ([]entity.EggPrice, error)
	GetEggPriceById(id uint64) (entity.EggPrice, error)
	UpdateEggPrice(eggPrice *entity.EggPrice) error
	DeleteEggPrice(id uint64) error

	CreateEggPriceDiscount(eggPriceDiscount *entity.EggPriceDiscount) error
	GetEggPriceDiscounts() ([]entity.EggPriceDiscount, error)
	GetEggPriceDiscountById(id uint64) (entity.EggPriceDiscount, error)
	UpdateEggPriceDiscount(eggPriceDiscount *entity.EggPriceDiscount) error
	DeleteEggPriceDiscount(id uint64) error
}

func NewEggPriceRepository(db *gorm.DB) IEggPriceRepository {
	return &EggPriceRepository{
		db: db,
	}
}

func (r *EggPriceRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *EggPriceRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *EggPriceRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *EggPriceRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *EggPriceRepository) CreateEggPrice(eggPrice *entity.EggPrice) error {
	return r.GetDB().Create(eggPrice).Error
}

func (r *EggPriceRepository) CreateEggPriceDiscount(eggPriceDiscount *entity.EggPriceDiscount) error {
	return r.GetDB().Create(eggPriceDiscount).Error
}

func (r *EggPriceRepository) GetEggPrices() ([]entity.EggPrice, error) {
	var eggPrice []entity.EggPrice
	err := r.GetDB().Find(&eggPrice).Error
	return eggPrice, err
}

func (r *EggPriceRepository) GetEggPriceDiscounts() ([]entity.EggPriceDiscount, error) {
	var eggPriceDiscount []entity.EggPriceDiscount
	err := r.GetDB().Find(&eggPriceDiscount).Error
	return eggPriceDiscount, err
}

func (r *EggPriceRepository) GetEggPriceById(id uint64) (entity.EggPrice, error) {
	var eggPrice entity.EggPrice
	err := r.GetDB().Where("id = ?", id).First(&eggPrice).Error
	return eggPrice, err
}

func (r *EggPriceRepository) GetEggPriceDiscountById(id uint64) (entity.EggPriceDiscount, error) {
	var eggPriceDiscount entity.EggPriceDiscount
	err := r.GetDB().Where("id = ?", id).First(&eggPriceDiscount).Error
	return eggPriceDiscount, err
}

func (r *EggPriceRepository) UpdateEggPrice(eggPrice *entity.EggPrice) error {
	return r.GetDB().Where("id = ?", eggPrice.Id).Updates(eggPrice).Error
}

func (r *EggPriceRepository) UpdateEggPriceDiscount(eggPriceDiscount *entity.EggPriceDiscount) error {
	return r.GetDB().Where("id = ?", eggPriceDiscount.Id).Updates(eggPriceDiscount).Error
}

func (r *EggPriceRepository) DeleteEggPrice(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.EggPrice{}).Error
}

func (r *EggPriceRepository) DeleteEggPriceDiscount(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.EggPriceDiscount{}).Error
}
