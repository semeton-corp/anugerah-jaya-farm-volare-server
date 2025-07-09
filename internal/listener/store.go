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

type StoreListener struct {
	cache cache.ICache
	db    *gorm.DB
	log   *zap.Logger
}

func NewStoreListener(cache cache.ICache, db *gorm.DB, log *zap.Logger) *StoreListener {
	return &StoreListener{cache: cache, db: db, log: log}
}

func (l *StoreListener) ListenStoreItemActivity(ctx context.Context, handler func(activity entity.StoreItemHistory)) {
	pubsub := l.cache.Subscribe(ctx, constant.StoreItemHistoryTopic)
	ch := pubsub.Channel()
	for msg := range ch {
		var activity entity.StoreItemHistory
		if err := json.Unmarshal([]byte(msg.Payload), &activity); err != nil {
			l.log.Error("failed to unmarshal store item activity", zap.Error(err))
			continue
		}
		l.log.Debug("sucess parsing data from redis", zap.Any("activity", activity))
		handler(activity)
	}
}

func (l *StoreListener) Start(ctx context.Context) {
	l.ListenStoreItemActivity(ctx, func(activity entity.StoreItemHistory) {
		err := l.db.Model(&entity.StoreItemHistory{}).Create(&activity).Error
		if err != nil {
			l.log.Error("failed to create store item history", zap.Error(err))
			return
		}
	})
}
