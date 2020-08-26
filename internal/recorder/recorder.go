package recorder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/team-bonitto/bonitto/internal/model"
	"github.com/team-bonitto/bonitto/internal/queue/consumer"
	"github.com/team-bonitto/bonitto/internal/queue/producer"
	"time"
)

const QueueName = "recorder"
const DBName = "db"

var _ consumer.Consumer = Recorder{}
var _ producer.Producer = Input{}

type Input model.Record

func (i Input) GetQueueName() string {
	return QueueName
}

func (i Input) Marshal() []byte {
	b, _ := json.Marshal(i)
	return b
}

type Recorder struct {
	Rds *redis.Client
}

func (r Recorder) GetQueueName() string {
	return QueueName
}

func (r Recorder) Consume(a string) error {
	input := Input{}
	if err := json.Unmarshal([]byte(a), &input); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	k := fmt.Sprintf("%s:%v", input.UserID, input.ProblemNo)
	if _, err := r.Rds.HMSet(ctx, DBName, k, a).Result(); err != nil {
		return err
	}
	return nil
}

func New(rds *redis.Client) (Recorder, error) {
	return Recorder{
		Rds: rds,
	}, nil
}

func (r *Recorder) GetResults(userID string, problemNo string) ([][]model.TestResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	k := fmt.Sprintf("%s:%v", userID, problemNo)
	a, err := r.Rds.HMGet(ctx, DBName, k).Result()
	if err != nil {
		return nil, err
	}
	if a[0] == nil {
		return nil, errors.New("not found " + k)
	}
	i := a[0].(string)
	input := Input{}
	if err := json.Unmarshal([]byte(i), &input); err != nil {
		return nil, err
	}
	return input.TestResults, nil
}
