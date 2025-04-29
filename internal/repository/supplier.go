package repository

import (
	"errors"

	"github.com/lib/pq"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type SupplierRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type ISupplierRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateSupplier(supplier *entity.Supplier) error
	GetSupplierById(id uint64) (entity.Supplier, error)
	GetAllSuppliers() ([]entity.Supplier, error)
	UpdateSupplier(supplier *entity.Supplier) error
	DeleteSupplier(id uint64) error
}

func NewSupplierRepository(db *gorm.DB) ISupplierRepository {
	return &SupplierRepository{
		db: db,
	}
}

func (r *SupplierRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *SupplierRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *SupplierRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *SupplierRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

// Todo : check to handle violation errror
func (r *SupplierRepository) CreateSupplier(supplier *entity.Supplier) error {
	err := r.GetDB().Create(&supplier).Error
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "foreign_key_violation" && pqErr.Constraint == "fk_suppliers_warehouse_item" {
				return errx.NotFound("Warehouse Item not found")
			}
		}

		return err
	}

	return nil
}

func (r *SupplierRepository) GetSupplierById(id uint64) (entity.Supplier, error) {
	var supplier entity.Supplier
	err := r.GetDB().Preload("WarehouseItem").Where("id = ?", id).First(&supplier).Error
	if err != nil {
		return entity.Supplier{}, err
	}

	return supplier, nil
}

func (r *SupplierRepository) GetAllSuppliers() ([]entity.Supplier, error) {
	var suppliers []entity.Supplier
	err := r.GetDB().Preload("WarehouseItem").Find(&suppliers).Error
	if err != nil {
		return nil, err
	}

	return suppliers, nil
}

func (r *SupplierRepository) UpdateSupplier(supplier *entity.Supplier) error {
	return r.GetDB().Model(&entity.Supplier{}).Where("id = ?", supplier.Id).Updates(supplier).Error
}

func (r *SupplierRepository) DeleteSupplier(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Supplier{}).Error
}
