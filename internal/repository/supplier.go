package repository

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
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

	CreateBatchSupplierItem(supplierItems *[]entity.SupplierItem) error
	DeleteSupplierItemInBatch(ids []uint64, supplierId uint64) error
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

func (r *SupplierRepository) CreateSupplier(supplier *entity.Supplier) error {
	err := r.GetDB().Create(&supplier).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *SupplierRepository) GetSupplierById(id uint64) (entity.Supplier, error) {
	var supplier entity.Supplier
	err := r.GetDB().Model(&entity.Supplier{}).Preload("SupplierItems.Item").Where("id = ?", id).First(&supplier).Error
	if err != nil {
		return entity.Supplier{}, err
	}

	return supplier, nil
}

func (r *SupplierRepository) GetAllSuppliers() ([]entity.Supplier, error) {
	var suppliers []entity.Supplier
	err := r.GetDB().Model(&entity.Supplier{}).Preload("SupplierItems.Item").Find(&suppliers).Error
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

// Todo : check to handle violation errror
func (r *SupplierRepository) CreateBatchSupplierItem(supplierItems *[]entity.SupplierItem) error {
	return r.GetDB().Model(&entity.SupplierItem{}).CreateInBatches(supplierItems, len(*supplierItems)).Error
}

func (r *SupplierRepository) DeleteSupplierItemInBatch(ids []uint64, supplierId uint64) error {
	return r.GetDB().Where("item_id IN ? AND supplier_id = ?", ids, supplierId).Delete(&entity.SupplierItem{}).Error
}
