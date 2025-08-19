package repository

import (
	"errors"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
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
	GetChickenMonitoringById(id uint64) (entity.ChickenMonitoring, error)
	UpdateChickenMonitoring(chickenMonitoring *entity.ChickenMonitoring) error
	GetChickenMonitorings(filter *dto.GetChickenMonitoringFilter) ([]entity.ChickenMonitoring, error)
	DeleteChickenMonitoring(id uint64) error

	CreateChickenHealthItem(chickenHealthItem *entity.ChickenHealthItem) error
	GetChickenHealthItems(filter dto.GetChickenHealthItemFilter) ([]entity.ChickenHealthItem, error)
	GetChickenHealthItemById(id uint64) (entity.ChickenHealthItem, error)
	UpdateChickenHealthItem(chickenHealthItem *entity.ChickenHealthItem) error
	DeleteChickenHealthItem(id uint64) error

	CreateChickenHealthMonitoring(chickenHealthMonitoring *entity.ChickenHealthMonitoring) error
	UpdateChickenHealthMonitoring(chickenHealthMonitoring *entity.ChickenHealthMonitoring) error
	GetChickenHealthMonitoringById(id uint64) (entity.ChickenHealthMonitoring, error)
	GetChickenHealthMonitoringByChickenCageId(chickenCageId uint64) ([]entity.ChickenHealthMonitoring, error)
	DeleteChickenHealthMonitoring(id uint64) error

	CountChickenMonitoringByChickenCageIdToday(cageId uint64) (int64, error)

	CreateChickenProcurementDraft(data *entity.ChickenProcurementDraft) error
	GetChickenProcurementDrafts() ([]entity.ChickenProcurementDraft, error)
	GetChickenProcurementDraft(id uint64) (entity.ChickenProcurementDraft, error)
	UpdateChickenProcurementDraft(data *entity.ChickenProcurementDraft) error
	DeleteChickenProcurementDraft(id uint64) error

	CreateChickenProcurement(data *entity.ChickenProcurement) error
	CreateChickenProcurementPaymentInBatch(data *[]entity.ChickenProcurementPayment) error
	GetChickenProcurements(filter dto.GetChickenProcurementFilter) ([]entity.ChickenProcurement, error)
	CountChickenProcurement(filter dto.GetChickenProcurementFilter) (int64, error)
	GetChickenProcurement(id uint64) (entity.ChickenProcurement, error)
	UpdateChickenProcurement(data *entity.ChickenProcurement) error

	CreateChickenProcurementPayment(data *entity.ChickenProcurementPayment) error
	GetChickenProcurementPayment(id uint64) (entity.ChickenProcurementPayment, error)
	UpdateChickenProcurementPayment(data *entity.ChickenProcurementPayment) error
	DeleteChickenProcurementPayment(id uint64) error

	CreateAfkirChickenCustomer(data *entity.AfkirChickenCustomer) error
	GetAfkirChickenCustomers() ([]entity.AfkirChickenCustomer, error)
	GetAfkirChickenCustomer(id uint64) (entity.AfkirChickenCustomer, error)
	UpdateAfkirChickenCustomer(data *entity.AfkirChickenCustomer) error
	DeleteAfkirChickenCustomer(id uint64) error

	CreateAfkirChickenSaleDraft(data *entity.AfkirChickenSaleDraft) error
	GetAfkirChickenSaleDrafts() ([]entity.AfkirChickenSaleDraft, error)
	GetAfkirChickenSaleDraft(id uint64) (entity.AfkirChickenSaleDraft, error)
	UpdateAfkirChickenSaleDraft(data *entity.AfkirChickenSaleDraft) error
	DeleteAfkirChickenSaleDraft(id uint64) error

	CreateAfkirChickenSale(data *entity.AfkirChickenSale) error
	CountChickenAfkirChickenSale(filter dto.GetAfkirChickenSaleFilter) (int64, error)
	GetAfkirChickenSales(filter dto.GetAfkirChickenSaleFilter) ([]entity.AfkirChickenSale, error)
	GetAfkirChickenSale(id uint64) (entity.AfkirChickenSale, error)
	UpdateAfkirChickenSale(data *entity.AfkirChickenSale) error

	DeleteAfkirChickenSalePayment(id uint64) error
	CreateAfkirChickenSalePayment(afkirChickenSalePayment *entity.AfkirChickenSalePayment) error
	GetAfkirChickenSalePaymentById(id uint64) (entity.AfkirChickenSalePayment, error)
	UpdateAfkirChickenSalePayment(afkirChickenSalePayment *entity.AfkirChickenSalePayment) error
	CreateAfkirChickenSalePaymentInBatch(data *[]entity.AfkirChickenSalePayment) error

	GetChickenPerformances(filter dto.GetChickenPerformanceFilter) ([]entity.ChickenPerformance, error)

	GetChickenMonitoringToday(chickenCageId uint64, date time.Time) (entity.ChickenMonitoring, error)
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

func (r *ChickenRepository) CountChickenMonitoringByChickenCageIdToday(cageId uint64) (int64, error) {
	var count int64
	if err := r.GetDB().Model(entity.ChickenMonitoring{}).Where("chicken_cage_id = ? AND DATE(created_at) = ?", cageId, time.Now()).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ChickenRepository) CreateChickenMonitoring(chickenMonitoring *entity.ChickenMonitoring) error {
	return r.GetDB().Create(chickenMonitoring).Error
}

func (r *ChickenRepository) GetChickenMonitoringById(id uint64) (entity.ChickenMonitoring, error) {
	var chickenMonitoring entity.ChickenMonitoring
	err := r.GetDB().
		Preload("ChickenCage.Cage.Location").
		Preload("ChickenCage.Cage.CagePlacement.User.Role").
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
	return r.GetDB().Model(&entity.ChickenMonitoring{}).Where("id = ?", chickenMonitoring.Id).Updates(map[string]interface{}{
		"chicken_cage_id":     chickenMonitoring.ChickenCageId,
		"total_death_chicken": chickenMonitoring.TotalDeathChicken,
		"total_sick_chicken":  chickenMonitoring.TotalSickChicken,
		"total_feed":          chickenMonitoring.TotalFeed,
		"note":                chickenMonitoring.Note,
		"update_by":           chickenMonitoring.UpdateBy,
		"updated_at":          chickenMonitoring.UpdatedAt,
	}).Error
}

func (r *ChickenRepository) GetChickenMonitorings(filter *dto.GetChickenMonitoringFilter) ([]entity.ChickenMonitoring, error) {
	var chickenMonitorings []entity.ChickenMonitoring
	query := r.GetDB().
		Preload("ChickenCage.Cage.Location").
		Preload("ChickenCage.Cage.CagePlacement.User.Role").
		Model(&entity.ChickenMonitoring{}).Joins("JOIN chicken_cages ON chicken_cages.id = chicken_monitorings.chicken_cage_id").Joins("JOIN cages ON cages.id = chicken_cages.cage_id")

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(chicken_monitorings.created_at) = ?", filter.Date.Value())
	}

	if filter.LocationId > 0 {
		query = query.Where("cages.location_id = ?", filter.LocationId)
	}

	if filter.CageId > 0 {
		query = query.Where("cages.id = ?", filter.CageId)
	}

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(chicken_monitorings.created_at) >= ? AND DATE(chicken_monitorings.created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	err := query.
		Order("created_at desc").
		Find(&chickenMonitorings).Error

	if err != nil {
		return nil, err
	}
	return chickenMonitorings, nil
}

func (r *ChickenRepository) DeleteChickenMonitoring(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ChickenMonitoring{}).Error
}

func (r *ChickenRepository) CreateChickenHealthItem(chickenHealthItem *entity.ChickenHealthItem) error {
	return r.GetDB().Model(&entity.ChickenHealthItem{}).Create(&chickenHealthItem).Error
}

func (r *ChickenRepository) GetChickenHealthItems(filter dto.GetChickenHealthItemFilter) ([]entity.ChickenHealthItem, error) {
	var chickenHealthItems []entity.ChickenHealthItem
	query := r.GetDB().Model(&entity.ChickenHealthItem{})

	if filter.Type.Value().IsValid() {
		query = query.Where("type = ?", filter.Type.Value())
	}

	err := query.Find(&chickenHealthItems).Error
	if err != nil {
		return nil, err
	}

	return chickenHealthItems, nil
}

func (r *ChickenRepository) GetChickenHealthItemById(id uint64) (entity.ChickenHealthItem, error) {
	var chickenHealthItem entity.ChickenHealthItem
	err := r.GetDB().Model(&entity.ChickenHealthItem{}).Where("id = ?", id).First(&chickenHealthItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ChickenHealthItem{}, errx.NotFound("chicken health item not found")
		}
		return entity.ChickenHealthItem{}, err
	}

	return chickenHealthItem, nil
}

func (r *ChickenRepository) UpdateChickenHealthItem(chickenHealthItem *entity.ChickenHealthItem) error {
	return r.GetDB().Model(&entity.ChickenHealthItem{}).Where("id = ?", chickenHealthItem.Id).Updates(&chickenHealthItem).Error
}

func (r *ChickenRepository) DeleteChickenHealthItem(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ChickenHealthItem{}).Error
}

func (r *ChickenRepository) CreateChickenHealthMonitoring(chickenHealthMonitoring *entity.ChickenHealthMonitoring) error {
	return r.GetDB().Model(&entity.ChickenHealthMonitoring{}).Create(&chickenHealthMonitoring).Error
}

func (r *ChickenRepository) UpdateChickenHealthMonitoring(chickenHealthMonitoring *entity.ChickenHealthMonitoring) error {
	return r.GetDB().Model(&entity.ChickenHealthMonitoring{}).Where("id = ?", chickenHealthMonitoring.Id).Updates(&chickenHealthMonitoring).Error
}

func (r *ChickenRepository) GetChickenHealthMonitoringById(id uint64) (entity.ChickenHealthMonitoring, error) {
	var chickenHealthMonitoring entity.ChickenHealthMonitoring

	err := r.GetDB().Model(&entity.ChickenHealthMonitoring{}).Where("id = ?", id).First(&chickenHealthMonitoring).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ChickenHealthMonitoring{}, errx.NotFound("chicken health monitoring not found")
		}
		return entity.ChickenHealthMonitoring{}, err
	}

	return chickenHealthMonitoring, nil
}

func (r *ChickenRepository) GetChickenHealthMonitoringByChickenCageId(chickenCageId uint64) ([]entity.ChickenHealthMonitoring, error) {
	chickenHealthMonitoring := make([]entity.ChickenHealthMonitoring, 0)
	err := r.GetDB().Model(&entity.ChickenHealthMonitoring{}).Find(&chickenHealthMonitoring).Error
	if err != nil {
		return chickenHealthMonitoring, err
	}

	return chickenHealthMonitoring, nil
}

func (r *ChickenRepository) DeleteChickenHealthMonitoring(id uint64) error {
	return r.GetDB().Model(&entity.ChickenHealthMonitoring{}).Where("id = ?", id).Delete(&entity.ChickenHealthMonitoring{}).Error
}

func (r *ChickenRepository) CreateChickenProcurementDraft(data *entity.ChickenProcurementDraft) error {
	return r.GetDB().Model(&entity.ChickenProcurementDraft{}).Create(data).Error
}

func (r *ChickenRepository) GetChickenProcurementDrafts() ([]entity.ChickenProcurementDraft, error) {
	chickenProcurementDrafts := make([]entity.ChickenProcurementDraft, 0)
	err := r.GetDB().Model(&entity.ChickenProcurementDraft{}).Preload("Cage.Location").Preload("Supplier").Order("total_price ASC").Find(&chickenProcurementDrafts).Error
	if err != nil {
		return nil, err
	}

	return chickenProcurementDrafts, nil
}

func (r *ChickenRepository) GetChickenProcurementDraft(id uint64) (entity.ChickenProcurementDraft, error) {
	var chickenProcurementDraft entity.ChickenProcurementDraft
	err := r.GetDB().Model(&entity.ChickenProcurementDraft{}).Where("id = ?", id).Preload("Cage.Location").Preload("Supplier").First(&chickenProcurementDraft).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ChickenProcurementDraft{}, errx.NotFound("chicken procurement draft not found")
		}
		return entity.ChickenProcurementDraft{}, err
	}

	return chickenProcurementDraft, nil
}

func (r *ChickenRepository) UpdateChickenProcurementDraft(data *entity.ChickenProcurementDraft) error {
	return r.GetDB().Model(&entity.ChickenProcurementDraft{}).Where("id = ?", data.Id).Updates(data).Error
}

func (r *ChickenRepository) DeleteChickenProcurementDraft(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ChickenProcurementDraft{}).Error
}

func (r *ChickenRepository) CreateChickenProcurement(data *entity.ChickenProcurement) error {
	return r.GetDB().Model(&entity.ChickenProcurement{}).Create(data).Error
}

func (r *ChickenRepository) CreateChickenProcurementPaymentInBatch(data *[]entity.ChickenProcurementPayment) error {
	return r.GetDB().Model(&entity.ChickenProcurementPayment{}).CreateInBatches(data, len(*data)).Error
}

func (r *ChickenRepository) UpdateChickenProcurement(data *entity.ChickenProcurement) error {
	return r.GetDB().Model(&entity.ChickenProcurement{}).Where("id = ?", data.Id).Updates(data).Error
}

func (r *ChickenRepository) GetChickenProcurement(id uint64) (entity.ChickenProcurement, error) {
	var chickenProcurement entity.ChickenProcurement
	err := r.GetDB().Model(&entity.ChickenProcurement{}).Preload("Cage.Location").Preload("Supplier").Preload("Payments").Where("id = ?", id).First(&chickenProcurement).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ChickenProcurement{}, errx.NotFound("chicken procurement not found")
		}
		return entity.ChickenProcurement{}, err
	}

	return chickenProcurement, err
}

func (r *ChickenRepository) CreateAfkirChickenSalePaymentInBatch(data *[]entity.AfkirChickenSalePayment) error {
	return r.GetDB().Model(&entity.AfkirChickenSalePayment{}).CreateInBatches(data, len(*data)).Error
}

func (r *ChickenRepository) CreateChickenProcurementPayment(data *entity.ChickenProcurementPayment) error {
	return r.GetDB().Model(&entity.ChickenProcurementPayment{}).Create(&data).Error
}

func (r *ChickenRepository) UpdateChickenProcurementPayment(data *entity.ChickenProcurementPayment) error {
	return r.GetDB().Model(&entity.ChickenProcurementPayment{}).Where("id = ?", data.Id).Updates(data).Error
}

func (r *ChickenRepository) GetChickenProcurementPayment(id uint64) (entity.ChickenProcurementPayment, error) {
	var chickenProcurementPayment entity.ChickenProcurementPayment
	err := r.GetDB().Model(&entity.ChickenProcurementPayment{}).Where("id = ?", id).First(&chickenProcurementPayment).Error
	if err != nil {
		return entity.ChickenProcurementPayment{}, err
	}

	return chickenProcurementPayment, nil
}

func (r *ChickenRepository) DeleteChickenProcurementPayment(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ChickenProcurementPayment{}).Error
}

func (r *ChickenRepository) CreateAfkirChickenCustomer(data *entity.AfkirChickenCustomer) error {
	err := r.GetDB().Model(&entity.AfkirChickenCustomer{}).Create(data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errx.BadRequest("afkir chicken customer already exist")
		}
		return err
	}

	return nil
}

func (r *ChickenRepository) GetAfkirChickenCustomers() ([]entity.AfkirChickenCustomer, error) {
	var data []entity.AfkirChickenCustomer
	err := r.GetDB().Model(&entity.AfkirChickenCustomer{}).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *ChickenRepository) GetAfkirChickenCustomer(id uint64) (entity.AfkirChickenCustomer, error) {
	var data entity.AfkirChickenCustomer
	err := r.GetDB().Model(&entity.AfkirChickenCustomer{}).Preload("AfkirChickenSales").Where("id = ?", id).First(&data).Error
	if err != nil {
		return entity.AfkirChickenCustomer{}, nil
	}

	return data, nil
}

func (r *ChickenRepository) UpdateAfkirChickenCustomer(data *entity.AfkirChickenCustomer) error {
	return r.GetDB().Model(&entity.AfkirChickenCustomer{}).Where("id = ?", data.Id).Updates(&data).Error
}

func (r *ChickenRepository) DeleteAfkirChickenCustomer(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.AfkirChickenCustomer{}).Error
}

func (r *ChickenRepository) CreateAfkirChickenSaleDraft(data *entity.AfkirChickenSaleDraft) error {
	return r.GetDB().Model(&entity.AfkirChickenSaleDraft{}).Create(data).Error
}

func (r *ChickenRepository) GetAfkirChickenSaleDrafts() ([]entity.AfkirChickenSaleDraft, error) {
	var data []entity.AfkirChickenSaleDraft
	err := r.GetDB().Model(&entity.AfkirChickenSaleDraft{}).Preload("ChickenCage.Cage.Location").Preload("ChickenCage.Cage.CagePlacement.User").Preload("ChickenCage.ChickenProcurement").Preload("AfkirChickenCustomer").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *ChickenRepository) GetAfkirChickenSaleDraft(id uint64) (entity.AfkirChickenSaleDraft, error) {
	var data entity.AfkirChickenSaleDraft
	err := r.GetDB().Model(&entity.AfkirChickenSaleDraft{}).Preload("ChickenCage.Cage.Location").Preload("ChickenCage.Cage.CagePlacement.User").Preload("ChickenCage.ChickenProcurement").Preload("AfkirChickenCustomer").Where("id = ?", id).First(&data).Error
	if err != nil {
		return entity.AfkirChickenSaleDraft{}, err
	}

	return data, nil
}

func (r *ChickenRepository) UpdateAfkirChickenSaleDraft(data *entity.AfkirChickenSaleDraft) error {
	return r.GetDB().Model(&entity.AfkirChickenSaleDraft{}).Where("id = ?", data.Id).Updates(data).Error
}

func (r *ChickenRepository) DeleteAfkirChickenSaleDraft(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.AfkirChickenSaleDraft{}).Error
}

func (r *ChickenRepository) CreateAfkirChickenSale(data *entity.AfkirChickenSale) error {
	return r.GetDB().Model(&entity.AfkirChickenSale{}).Create(data).Error
}

func (r *ChickenRepository) GetAfkirChickenSales(filter dto.GetAfkirChickenSaleFilter) ([]entity.AfkirChickenSale, error) {
	var data []entity.AfkirChickenSale
	query := r.GetDB().Model(&entity.AfkirChickenSale{})

	if filter.Page > 0 {
		query = query.Limit(int(constant.PaginationDefaultLimit)).Offset((int(filter.Page) - 1) * int(constant.PaginationDefaultLimit))
	}

	err := query.Preload("ChickenCage.Cage.Location").Preload("ChickenCage.Cage.CagePlacement.User").Preload("ChickenCage.ChickenProcurement").Preload("AfkirChickenCustomer").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *ChickenRepository) CountChickenAfkirChickenSale(filter dto.GetAfkirChickenSaleFilter) (int64, error) {
	var count int64
	err := r.GetDB().Model(entity.AfkirChickenSale{}).Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *ChickenRepository) GetAfkirChickenSale(id uint64) (entity.AfkirChickenSale, error) {
	var data entity.AfkirChickenSale
	err := r.GetDB().Model(&entity.AfkirChickenSale{}).Where("id = ?", id).Preload("ChickenCage.Cage.Location").Preload("ChickenCage.Cage.CagePlacement.User").Preload("ChickenCage.ChickenProcurement").Preload("AfkirChickenCustomer").Preload("Payments").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.AfkirChickenSale{}, errx.NotFound("afkir chicken sale not found")
		}
		return entity.AfkirChickenSale{}, nil
	}

	return data, nil
}

func (r *ChickenRepository) DeleteAfkirChickenSalePayment(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.AfkirChickenSalePayment{}).Error
}

func (r *ChickenRepository) CreateAfkirChickenSalePayment(afkirChickenSalePayment *entity.AfkirChickenSalePayment) error {
	return r.GetDB().Model(&entity.AfkirChickenSalePayment{}).Create(afkirChickenSalePayment).Error
}

func (r *ChickenRepository) GetAfkirChickenSalePaymentById(id uint64) (entity.AfkirChickenSalePayment, error) {
	var afkirChickenSalePayment entity.AfkirChickenSalePayment
	err := r.GetDB().Model(&entity.AfkirChickenSalePayment{}).First(&afkirChickenSalePayment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.AfkirChickenSalePayment{}, errx.NotFound("afkir chicken sale payment not found")
		}
		return entity.AfkirChickenSalePayment{}, err
	}
	return afkirChickenSalePayment, nil
}

func (r *ChickenRepository) UpdateAfkirChickenSalePayment(afkirChickenSalePayment *entity.AfkirChickenSalePayment) error {
	return r.GetDB().Model(entity.AfkirChickenSalePayment{}).Where("id = ?", afkirChickenSalePayment.Id).Updates(afkirChickenSalePayment).Error
}

func (r *ChickenRepository) UpdateAfkirChickenSale(data *entity.AfkirChickenSale) error {
	return r.GetDB().Model(&entity.AfkirChickenSale{}).Where("id = ?", data.Id).Updates(data).Error
}

func (r *ChickenRepository) GetChickenProcurements(filter dto.GetChickenProcurementFilter) ([]entity.ChickenProcurement, error) {
	var chickenProcurements []entity.ChickenProcurement
	query := r.GetDB().Model(&entity.ChickenProcurement{})

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	if filter.Page > 0 {
		query = query.Limit(int(constant.PaginationDefaultLimit)).Offset((int(filter.Page) - 1) * int(constant.PaginationDefaultLimit))
	}

	err := query.Preload("Cage.Location").Preload("Supplier").Find(&chickenProcurements).Error
	if err != nil {
		return nil, err
	}

	return chickenProcurements, nil
}

func (r *ChickenRepository) CountChickenProcurement(filter dto.GetChickenProcurementFilter) (int64, error) {
	var count int64
	query := r.GetDB().Model(&entity.ChickenProcurement{})

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	err := query.Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *ChickenRepository) GetChickenPerformances(filter dto.GetChickenPerformanceFilter) ([]entity.ChickenPerformance, error) {
	var chickenPerformances []entity.ChickenPerformance
	query := r.GetDB().Model(&entity.ChickenPerformance{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.LocationId > 0 {
		query = query.Where("location_id = ?", filter.LocationId)
	}

	err := query.Find(&chickenPerformances).Error
	if err != nil {
		return nil, err
	}

	return chickenPerformances, nil
}

func (r *ChickenRepository) GetChickenMonitoringToday(chickenCageId uint64, date time.Time) (entity.ChickenMonitoring, error) {
	var monitoring entity.ChickenMonitoring

	err := r.GetDB().
		Where("chicken_cage_id = ? AND DATE(created_at) = ?", chickenCageId, date.Format("2006-01-02")).
		First(&monitoring).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ChickenMonitoring{}, nil
		}
		return entity.ChickenMonitoring{}, err
	}

	return monitoring, nil
}
