package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/logger"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, here interface{}) error
	Del(ctx context.Context, key string) error
	Alive(ctx context.Context) bool
}

type redisCache struct {
	client *redis.Client
	ttl    time.Duration
	logger logger.Logger
}

func NewRedisCache(logger logger.Logger, ttl time.Duration, client *redis.Client) Cache {
	return &redisCache{
		client: client,
		ttl:    ttl,
		logger: logger,
	}
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}) error {
	log := c.logger.WithField("key", key).WithField("value", value)
	b, err := json.Marshal(value)
	if err != nil {
		c.logger.WithError(err).Error(ctx, "cache_error")
		return err
	}
	log.Info(ctx, "Saving Value to Key")
	err = c.client.Set(context.Background(), key, string(b), c.ttl).Err()
	if err != nil {
		log.WithError(err).Error(ctx, "cache_error")
		return err
	}
	return nil
}

func (c *redisCache) Get(ctx context.Context, key string, here interface{}) error {
	log := c.logger.WithField("key", key)

	log.Info(ctx, "Retrieving Key")
	val, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		log.WithError(err).Error(ctx, "cache_error")
		return err
	}
	err = json.Unmarshal([]byte(val), here)
	if err != nil {
		log.WithError(err).Error(ctx, "cache_error")
		return err
	}
	return nil
}

func (c *redisCache) Del(ctx context.Context, key string) error {
	log := c.logger.WithField("key", key)

	log.WithField("key", key).Info(ctx, "Deleting Key")
	numErased, err := c.client.Del(context.Background(), key).Result()
	if err != nil {
		log.WithError(err).Error(ctx, "cache_error")
		return err
	}
	if numErased == 0 {
		log.Error(ctx, "cache key not found")
		return redis.Nil
	}

	return nil
}

func (c *redisCache) Alive(ctx context.Context) bool {
	c.logger.Info(ctx, "Pinging Redis")
	err := c.client.Ping(context.Background()).Err()
	if err != nil {
		c.logger.WithError(err).Error(ctx, "cache not connected")
		return false
	}
	return true
}
