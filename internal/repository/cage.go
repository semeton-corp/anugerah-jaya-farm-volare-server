package repository

import (
	"errors"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
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
	GetCagesByIds(ids []uint64) ([]entity.Cage, error)

	CreateChickenCage(chickenCage *entity.ChickenCage) error
	UpdateChickenCage(chickenCage *entity.ChickenCage) error
	CreateChickenCageInBatch(chickenCage *[]entity.ChickenCage) error
	GetChickenCages(filter dto.GetChickenCageFilter) ([]entity.ChickenCage, error)
	GetChickenCageById(id uint64) (entity.ChickenCage, error)
	GetChickenCagesByCageIds(ids []uint64) ([]entity.ChickenCage, error)
	GetChickenCageByCageId(cageId uint64) (entity.ChickenCage, error)
	GetChickenCageByIds(ids []uint64) ([]entity.ChickenCage, error)

	CreateCageFeed(data *entity.CageFeed) error
	GetCageFeeds() ([]entity.CageFeed, error)
	UpdateCageFeed(data *entity.CageFeed) error
	GetCageFeed(id uint64) (entity.CageFeed, error)
	GetCageFeedByChickenCategory(chickenCategory enum.ChickenCategory) (entity.CageFeed, error)
	CreateCageFeedDetail(data *entity.CageFeedDetail) error
	UpdateCageFeedDetail(data *entity.CageFeedDetail) error
	CreateCageFeedDetails(details *[]entity.CageFeedDetail) error
	DeleteCageFeedDetailsNotIn(cageFeedId uint64, ids []uint64) error
	UpsertCageFeedDetails(details *[]entity.CageFeedDetail) error

	GetChickenCageFeeds(filter dto.GetChickenCageFeedFilter) ([]entity.ChickenCage, error)

	CreateCageFeedStock(data *entity.CageFeedStock) error
	UpdateCageFeedStock(data *entity.CageFeedStock) error
	GetCageFeedStock(id uint64) (entity.CageFeedStock, error)
	GetCageFeedStocks(filter dto.GetCageFeedStockFilter) ([]entity.CageFeedStock, error)
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
		query = query.Where("location_id = ?", filter.LocationId)
	}

	if filter.ChickenCategory.Value().IsValid() {
		query = query.Where("chicken_category = ?", filter.ChickenCategory.Value())
	}

	if filter.IsUsed != nil {
		query = query.Where("is_used = ?", filter.ChickenCategory)
	}

	err := query.Order("created_at DESC").Preload("Location").Find(&cages).Error
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
	return r.GetDB().Model(&entity.Cage{}).Where("id = ?", cage.Id).Updates(map[string]interface{}{
		"location_id":      cage.LocationId,
		"name":             cage.Name,
		"chicken_category": cage.ChickenCategory,
		"is_used":          cage.IsUsed,
		"updated_by":       cage.UpdatedBy,
	}).Error
}

func (r *CageRepository) DeleteCage(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Cage{}).Error
}

func (r *CageRepository) GetChickenCageByCageId(cageId uint64) (entity.ChickenCage, error) {
	var chickenCage entity.ChickenCage
	// Note : find the newest chicken cage by cage id, that's why sort desc
	err := r.GetDB().Model(&entity.ChickenCage{}).Where("cage_id = ?", cageId).Order("created_at DESC").First(&chickenCage).Error

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

	if filter.CageId > 0 {
		query = query.Where("chicken_cages.cage_id = ?", filter.CageId)
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

func (r *CageRepository) GetCagesByIds(ids []uint64) ([]entity.Cage, error) {
	var cages []entity.Cage
	err := r.GetDB().Model(&entity.Cage{}).Where("id IN ?", ids).Order("created_at DESC").Find(&cages).Error
	if err != nil {
		return nil, err
	}

	return cages, nil
}

func (r *CageRepository) GetChickenCagesByCageIds(cageIds []uint64) ([]entity.ChickenCage, error) {
	var chickenCages []entity.ChickenCage

	subQuery := r.GetDB().
		Model(&entity.ChickenCage{}).
		Select("DISTINCT ON (cage_id) *").
		Where("cage_id IN ?", cageIds).
		Order("cage_id, created_at DESC")

	err := r.GetDB().
		Table("(?) as chicken_cages", subQuery).
		Preload("Cage.CagePlacement.User.Role").
		Find(&chickenCages).Error

	if err != nil {
		return nil, err
	}

	return chickenCages, nil
}

func (r *CageRepository) GetChickenCageByIds(ids []uint64) ([]entity.ChickenCage, error) {
	var chickenCages []entity.ChickenCage

	err := r.GetDB().Model(entity.ChickenCage{}).Where("id IN ?", ids).Preload("ChickenProcurement").Preload("Cage.Location").Preload("Cage.CagePlacement.User.Role").Find(&chickenCages).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}

	return chickenCages, nil
}

func (r *CageRepository) UpdateChickenCage(chickenCage *entity.ChickenCage) error {
	return r.GetDB().
		Model(&entity.ChickenCage{}).
		Where("id = ?", chickenCage.Id).
		Updates(map[string]interface{}{
			"cage_id":                            chickenCage.CageId,
			"chicken_procurement_id":             chickenCage.ChickenProcurementId,
			"total_chicken":                      chickenCage.TotalChicken,
			"is_need_routine_vaccine":            chickenCage.IsNeedRoutineVaccine,
			"latest_chicken_age_vaccine_routine": chickenCage.LatestChickenAgeVaccineRoutine,
			"is_need_feed":                       chickenCage.IsNeedFeed,
			"updated_at":                         time.Now(),
			"updated_by":                         chickenCage.UpdatedBy,
		}).Error
}

func (r *CageRepository) CreateChickenCageInBatch(chickenCage *[]entity.ChickenCage) error {
	return r.GetDB().Model(&entity.ChickenCage{}).CreateInBatches(&chickenCage, len(*chickenCage)).Error
}

func (r *CageRepository) CreateCageFeed(data *entity.CageFeed) error {
	return r.GetDB().Model(&entity.CageFeed{}).Create(&data).Error
}

func (r *CageRepository) GetCageFeeds() ([]entity.CageFeed, error) {
	var cageFeeds []entity.CageFeed
	err := r.GetDB().Model(&entity.CageFeed{}).Preload("CageFeedDetails.Item").Find(&cageFeeds).Error
	if err != nil {
		return nil, err
	}

	return cageFeeds, nil
}

func (r *CageRepository) UpdateCageFeed(data *entity.CageFeed) error {
	return r.GetDB().Model(&entity.CageFeed{}).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"chicken_category": data.ChickenCategory,
		"feed_type":        data.FeedType,
		"total_feed":       data.TotalFeed,
		"updated_by":       data.UpdatedBy,
	}).Error
}

func (r *CageRepository) GetCageFeed(id uint64) (entity.CageFeed, error) {
	var cageFeed entity.CageFeed
	err := r.GetDB().Model(&entity.CageFeed{}).Where("id = ?", id).Preload("CageFeedDetails.Item").First(&cageFeed).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.CageFeed{}, errx.NotFound("cage feed not found")
		}
		return entity.CageFeed{}, err
	}

	return cageFeed, nil
}

func (r *CageRepository) CreateCageFeedDetail(data *entity.CageFeedDetail) error {
	return r.GetDB().Model(&entity.CageFeedDetail{}).Create(&data).Error
}

func (r *CageRepository) UpdateCageFeedDetail(data *entity.CageFeedDetail) error {
	return r.GetDB().Model(&entity.CageFeedDetail{}).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"cage_feed_id": data.CageFeedId,
		"item_id":      data.ItemId,
		"percentage":   data.Percentage,
		"updated_by":   data.UpdatedBy,
		"updated_at":   time.Now(),
	}).Error
}

func (r *CageRepository) CreateCageFeedDetails(details *[]entity.CageFeedDetail) error {
	return r.GetDB().Model(&entity.CageFeedDetail{}).Create(details).Error
}

func (r *CageRepository) DeleteCageFeedDetailsNotIn(cageFeedId uint64, ids []uint64) error {
	query := r.GetDB().Where("cage_feed_id = ?", cageFeedId)
	if len(ids) > 0 {
		query = query.Where("id NOT IN ?", ids)
	}
	return query.Delete(&entity.CageFeedDetail{}).Error
}

func (r *CageRepository) UpsertCageFeedDetails(details *[]entity.CageFeedDetail) error {
	for _, detail := range *details {
		var existing entity.CageFeedDetail
		err := r.GetDB().Where("id = ?", detail.Id).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := r.GetDB().Create(&detail).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			if err := r.GetDB().Model(&entity.CageFeedDetail{}).Where("id = ?", detail.Id).Updates(&detail).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *CageRepository) GetChickenCageFeeds(filter dto.GetChickenCageFeedFilter) ([]entity.ChickenCage, error) {
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
		Order("chicken_cages.created_at DESC").
		Find(&chickenCages).Error

	if err != nil {
		return nil, err
	}

	return chickenCages, nil
}

func (r *CageRepository) CreateCageFeedStock(data *entity.CageFeedStock) error {
	return r.GetDB().Model(&entity.CageFeedStock{}).Create(data).Error
}

func (r *CageRepository) UpdateCageFeedStock(data *entity.CageFeedStock) error {
	return r.GetDB().Model(&entity.CageFeedStock{}).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"cage_id":    data.CageId,
		"total_feed": data.TotalFeed,
		"used_feed":  data.UsedFeed,
		"updated_by": data.UpdatedBy,
	}).Error
}

func (r *CageRepository) GetCageFeedStock(id uint64) (entity.CageFeedStock, error) {
	var data entity.CageFeedStock
	err := r.GetDB().Model(&entity.CageFeedStock{}).Where("id = ?", id).First(&data).Error
	if err != nil {
		return entity.CageFeedStock{}, err
	}

	return data, nil
}

func (r *CageRepository) GetCageFeedStocks(filter dto.GetCageFeedStockFilter) ([]entity.CageFeedStock, error) {
	var data []entity.CageFeedStock
	query := r.GetDB().Model(&entity.CageFeedStock{})

	if filter.CageId > 0 {
		query = query.Where("cage_id = ?", filter.CageId)
	}

	err := query.Where("total_feed - used_feed <> 0").Order("created_by ASC").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *CageRepository) GetCageFeedByChickenCategory(chickenCategory enum.ChickenCategory) (entity.CageFeed, error) {
	var cageFeed entity.CageFeed
	err := r.GetDB().Model(&entity.CageFeed{}).Where("chicken_category = ?", chickenCategory).Preload("CageFeedDetails.Item").First(&cageFeed).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.CageFeed{}, errx.NotFound("cage feed not found")
		}
		return entity.CageFeed{}, err
	}

	return cageFeed, nil
}
