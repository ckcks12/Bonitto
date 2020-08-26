package queue

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewRedisClient(o *redis.Options) (*redis.Client, error){
	rdb := redis.NewClient(o)
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return rdb, nil
}
