package listener

import (
	"context"
	"encoding/json"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/cache"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotificationListener struct {
	cache cache.ICache
	db    *gorm.DB
	log   *zap.Logger
}

func NewNotificationListener(cache cache.ICache, db *gorm.DB, log *zap.Logger) *NotificationListener {
	return &NotificationListener{cache: cache, db: db, log: log}
}

func (l *NotificationListener) ListenNotificationItemActivity(ctx context.Context, handler func(notification entity.Notification)) {
	pubsub := l.cache.Subscribe(ctx, constant.NotificationTopic)
	ch := pubsub.Channel()
	for msg := range ch {
		var activity entity.Notification
		if err := json.Unmarshal([]byte(msg.Payload), &activity); err != nil {
			l.log.Error("failed to unmarshal notification", zap.Error(err))
			continue
		}
		l.log.Debug("success parsing data from redis", zap.Any("notification", activity))
		handler(activity)
	}
}

func (l *NotificationListener) Start(ctx context.Context) {
	l.ListenNotificationItemActivity(ctx, func(activity entity.Notification) {
		err := l.db.Model(&entity.Notification{}).Create(&activity).Error
		if err != nil {
			l.log.Error("failed to create notification", zap.Error(err))
			return
		}
	})
}
