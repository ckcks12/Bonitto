package producer

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/team-bonitto/bonitto/internal/queue"
)

type Producer interface {
	GetQueueName() string
	Marshal() []byte
}

type RedisProducer struct {
	RDB *redis.Client
}

func New(addr string) (*RedisProducer, error) {
	rdb, err := queue.NewRedisClient(&redis.Options{Addr: addr})
	if err != nil {
		return nil, err
	}
	p := &RedisProducer{RDB: rdb}
	return p, nil
}

func (p *RedisProducer) Produce(a Producer) error {
	n := a.GetQueueName()
	b := a.Marshal()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := p.RDB.RPush(ctx, n, b).Result(); err != nil {
		return err
	}
	return nil
}
