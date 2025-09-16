package repository

import (
	"strings"

	"github.com/lib/pq"
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
	db := r.GetDB().Model(&entity.Notification{})

	var ors []string
	var args []interface{}

	if filter.StoreIds != nil {
		if len(filter.NotificationContexts) > 0 {
			ors = append(ors, "(store_id IN ? AND notification_contexts && ?)")
			args = append(args, filter.StoreIds, pq.StringArray(filter.NotificationContexts))
		} else {
			ors = append(ors, "store_id IN ?")
			args = append(args, filter.StoreIds)
		}
	}

	if filter.CageIds != nil {
		if len(filter.NotificationContexts) > 0 {
			ors = append(ors, "(cage_id IN ? AND notification_contexts && ?)")
			args = append(args, filter.CageIds, pq.StringArray(filter.NotificationContexts))
		} else {
			ors = append(ors, "cage_id IN ?")
			args = append(args, filter.CageIds)
		}
	}

	if filter.WarehouseIds != nil {
		if len(filter.NotificationContexts) > 0 {
			ors = append(ors, "(warehouse_id IN ? AND notification_contexts && ?)")
			args = append(args, filter.WarehouseIds, pq.StringArray(filter.NotificationContexts))
		} else {
			ors = append(ors, "warehouse_id IN ?")
			args = append(args, filter.WarehouseIds)
		}
	}

	if len(ors) > 0 {
		db = db.Where(strings.Join(ors, " OR "), args...)
	}

	if filter.UserIds != nil {
		if len(ors) > 0 || filter.IsMarked != nil {
			sub := r.GetDB().Model(&entity.Notification{})
			sub = sub.Where(strings.Join(ors, " OR "), args...)
			if filter.IsMarked != nil {
				sub = sub.Where("is_marked = ?", filter.IsMarked)
			}
			db = r.GetDB().Model(&entity.Notification{}).
				Where("user_id IN ? OR id IN (?)", filter.UserIds, sub.Select("id"))
		} else {
			db = db.Where("user_id IN ?", filter.UserIds)
		}
	}

	if filter.IsMarked != nil {
		db = db.Where("is_marked = ?", filter.IsMarked)
	}

	err := db.Order("created_at DESC").Find(&data).Error
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
	err := r.GetDB().Order("created_at DESC").Model(&entity.Notification{}).Where("id IN ?", ids).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}
