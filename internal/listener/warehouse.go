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

type WarehouseListener struct {
	cache cache.ICache
	db    *gorm.DB
	log   *zap.Logger
}

func NewWarehouseListener(cache cache.ICache, db *gorm.DB, log *zap.Logger) *WarehouseListener {
	return &WarehouseListener{cache: cache, db: db, log: log}
}

func (l *WarehouseListener) ListenWarehouseItemActivity(ctx context.Context, handler func(activity entity.WarehouseItemHistory)) {
	pubsub := l.cache.Subscribe(ctx, constant.WarehouseItemHistoryTopic)
	ch := pubsub.Channel()
	for msg := range ch {
		var activity entity.WarehouseItemHistory
		if err := json.Unmarshal([]byte(msg.Payload), &activity); err != nil {
			l.log.Error("failed to unmarshal warehouse item activity", zap.Error(err))
			continue
		}
		l.log.Debug("sucess parsing data from redis", zap.Any("activity", activity))
		handler(activity)
	}
}

func (l *WarehouseListener) Start(ctx context.Context) {
	l.ListenWarehouseItemActivity(ctx, func(activity entity.WarehouseItemHistory) {
		err := l.db.Model(&entity.WarehouseItemHistory{}).Create(&activity).Error
		if err != nil {
			l.log.Error("failed to create warehouse item history", zap.Error(err))
			return
		}
	})
}
