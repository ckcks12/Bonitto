package consumer

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/team-bonitto/bonitto/internal/queue"
)

type Consumer interface {
	GetQueueName() string
	Consume(a string) error
}

type RedisConsumer struct {
	RDB *redis.Client
}

func New(addr string) (*RedisConsumer, error) {
	rdb, err := queue.NewRedisClient(&redis.Options{Addr: addr})
	if err != nil {
		return nil, err
	}
	p := &RedisConsumer{RDB: rdb}
	return p, nil
}

func (s *RedisConsumer) Consume(c Consumer) error {
	str, err := s.ConsumeManually(c.GetQueueName())
	if err != nil {
		return err
	}
	if err := c.Consume(str); err != nil {
		return err
	}
	return nil
}

func (s *RedisConsumer) ConsumeManually(queueName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	str, err := s.RDB.LPop(ctx, queueName).Result()
	if err != nil {
		return "", err
	}
	return str, nil
}
