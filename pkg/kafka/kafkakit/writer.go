package kafkakit

import (
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	ErrBrokersRequired = errors.New("brokers required")
)

type Writer struct {
	brokers []string
	topic   string

	*kafka.Writer
}

// NewWriter creates a new instance of kafka.Writer with the specified parameters
func NewWriter(brokers []string, topic string, compress bool, batchSize int, batchTimeout time.Duration,
	requiredAcks string, createTopic bool, maxReqSize int64) (*Writer, error) {

	if len(brokers) == 0 {
		return nil, ErrBrokersRequired
	}

	if batchSize <= 0 {
		batchSize = 50
	}

	if batchTimeout <= 0 {
		batchTimeout = 10
	}

	if maxReqSize <= 0 {
		maxReqSize = 1048576
	}

	if requiredAcks == "" {
		requiredAcks = "one"
	}

	var acks kafka.RequiredAcks
	err := acks.UnmarshalText([]byte(requiredAcks))
	if err != nil {
		return nil, err
	}

	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Balancer:               &kafka.Hash{},
		Topic:                  topic,
		RequiredAcks:           acks,
		BatchSize:              batchSize,
		BatchTimeout:           batchTimeout,
		AllowAutoTopicCreation: createTopic,
		BatchBytes:             maxReqSize,
	}

	if compress {
		kafkaWriter.Compression = kafka.Zstd
	}

	return &Writer{
		Writer:  kafkaWriter,
		brokers: brokers,
		topic:   topic,
	}, nil
}
