package broker

import (
	"context"
	"fmt"
	"github.com/eugene-vodyanko/orgonaut/internal/model"
	"github.com/eugene-vodyanko/orgonaut/pkg/kafka/kafkakit"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

type Broker struct {
	writer *kafkakit.Writer
}

func NewBroker(writer *kafkakit.Writer) *Broker {
	return &Broker{
		writer,
	}
}

// SendRecords sends messages to Kafka in the specified topic.
// A text representation is used as the key of the Kafka message (e.g., "id=32").
// A flat JSON representation is used as the value of the Kafka message (e.g., {"col_name":"col_value", ...})
func (b *Broker) SendRecords(ctx context.Context, topic string, records []*model.Record) error {
	start := time.Now()

	kafkaMessages := make([]kafka.Message, 0, len(records))
	for _, record := range records {

		key, err := record.GetKey()
		if err != nil {
			return fmt.Errorf("broker - get key failed: %w", err)
		}

		value, err := record.GetValue()
		if err != nil {
			return fmt.Errorf("broker - get value failed: %w", err)
		}

		kafkaMessages = append(kafkaMessages,
			kafka.Message{
				Key:   key,
				Value: value,
				Topic: topic,
			},
		)
	}

	err := b.writer.WriteMessages(ctx, kafkaMessages...)
	if err != nil {
		return fmt.Errorf("broker - write messages failed: %w", err)
	}

	elapsed := time.Now()

	slog.Debug("broker - write to kafka",
		"elapsed", elapsed.Sub(start),
		"topic", topic,
	)

	return nil
}
