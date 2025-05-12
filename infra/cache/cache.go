package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/spf13/viper"
)

type ICache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
	Publish(ctx context.Context, topic string, value string) error
	Subscribe(ctx context.Context, topic string) *redis.PubSub
}

type Cache struct {
	client *redis.Client
}

func New() ICache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("cache.host"), viper.GetInt("cache.port")),
		Password: viper.GetString("cache.password"),
		DB:       viper.GetInt("cache.db"),
	})

	return &Cache{client: rdb}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	res := c.client.Get(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		return "", nil
	}
	return res.Result()
}

func (c *Cache) Set(ctx context.Context, key string, value string) error {
	return c.client.Set(context.Background(), key, value, time.Duration(constant.CacheDefaultDuration)).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Publish(ctx context.Context, topic string, value string) error {
	return c.client.Publish(ctx, topic, value).Err()
}

func (c *Cache) Subscribe(ctx context.Context, topic string) *redis.PubSub {
	subscriber := c.client.Subscribe(ctx, topic)
	return subscriber
}
