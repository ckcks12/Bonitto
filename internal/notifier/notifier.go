package notifier

import (
	"encoding/json"
	"github.com/team-bonitto/bonitto/internal/model"
	"github.com/team-bonitto/bonitto/internal/queue/consumer"
	"github.com/team-bonitto/bonitto/internal/queue/producer"
)

const QueueName = "notifier"

var _ producer.Producer = Input{}
var _ consumer.Consumer = Notifier{}

type Notifier struct{}

type Input model.Notification

func (i Input) GetQueueName() string {
	return QueueName
}

func (i Input) Marshal() []byte {
	b, _ := json.Marshal(i)
	return b
}

func (n Notifier) GetQueueName() string {
	return QueueName
}

func (n Notifier) Consume(a string) error {
	input := Input{}
	if err := json.Unmarshal([]byte(a), &input); err != nil {
		return err
	}
	// TODO: notify logic
	return nil
}

func New() (Notifier, error) {
	return Notifier{}, nil
}
