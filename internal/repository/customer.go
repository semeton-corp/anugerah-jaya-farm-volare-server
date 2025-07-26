package repository

import (
	"errors"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type ICustomerRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetCustomers() ([]entity.Customer, error)
	GetCustomerById(id uint64) (entity.Customer, error)
	CreateCustomer(data *entity.Customer) error
	DeleteCustomer(id uint64) error
}

func NewCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *CustomerRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *CustomerRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *CustomerRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *CustomerRepository) GetCustomers() ([]entity.Customer, error) {
	var customers []entity.Customer
	err := r.GetDB().Model(&entity.Customer{}).Preload("StoreSales").Preload("WarehouseSales").Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) CreateCustomer(data *entity.Customer) error {
	err := r.GetDB().Model(entity.Customer{}).Create(data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errx.BadRequest("customer already exist")
		}
		return err
	}

	return nil
}

func (r *CustomerRepository) GetCustomerById(id uint64) (entity.Customer, error) {
	var customer entity.Customer
	err := r.GetDB().Model(&entity.Customer{}).Where("id = ?", id).Preload("StoreSales").Preload("WarehouseSales").First(&customer).Error
	if err != nil {
		return entity.Customer{}, err
	}

	return customer, nil
}

func (r *CustomerRepository) DeleteCustomer(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Customer{}).Error
}
