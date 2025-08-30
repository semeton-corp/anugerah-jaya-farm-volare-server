package repository

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type INotificationRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateNotification(data *entity.Notification) error
	GetNotifications(filter dto.GetNotificationFilter) ([]entity.Notification, error)
	UpdateNotificationAsMarked(ids []uint64) error
	GetNotificationByIds(ids []uint64) ([]entity.Notification, error)
}

func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *NotificationRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *NotificationRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *NotificationRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *NotificationRepository) CreateNotification(data *entity.Notification) error {
	return r.GetDB().Model(&entity.Notification{}).Create(data).Error
}

func (r *NotificationRepository) GetNotifications(filter dto.GetNotificationFilter) ([]entity.Notification, error) {
	var data []entity.Notification
	query := r.GetDB().Model(&entity.Notification{})

	if filter.CageId > 0 {
		query = query.Where("cage_id = ?", filter.CageId)
	}

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	if filter.StoreId > 0 {
		query = query.Where("store_id = ?", filter.StoreId)
	}

	if filter.UserId != "" {
		query = query.Where("user_id = ?", filter.UserId)
	}

	if filter.IsMarked != nil {
		query = query.Where("is_marked = ?", filter.IsMarked)
	}

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *NotificationRepository) UpdateNotificationAsMarked(ids []uint64) error {
	return r.GetDB().Model(&entity.Notification{}).Where("id IN ?", ids).Updates(map[string]any{
		"is_marked": true,
	}).Error
}

func (r *NotificationRepository) GetNotificationByIds(ids []uint64) ([]entity.Notification, error) {
	var data []entity.Notification
	err := r.GetDB().Model(&entity.Notification{}).Where("id IN ?", ids).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}
