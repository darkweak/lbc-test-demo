package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"leboncoin/pkg/services/pubsub"

	"github.com/segmentio/kafka-go"
)

const maxMessageBytes = 10 * 1024 * 1024 // 10 MiB

type kafkaConsumer struct {
	reader *kafka.Reader
}

var _ pubsub.Consumer = (*kafkaConsumer)(nil)

func NewKafkaConsumer(hosts []string, topic string) *kafkaConsumer {
	return &kafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  hosts,
			GroupID:  "consumer-" + topic,
			Topic:    topic,
			MaxBytes: maxMessageBytes,
		}),
	}
}

func (k kafkaConsumer) Consume(callback func(message pubsub.Message) error) error {
	ctx := context.Background()
	for {
		kafkaMsg, err := k.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("read kafka message: %w", err)
		}

		var message pubsub.Message

		err = json.Unmarshal(kafkaMsg.Value, &message)
		if err != nil {
			return fmt.Errorf("unmarshal kafka message: %w", err)
		}

		err = callback(message)
		if err != nil {
			return fmt.Errorf("callback kafka message: %w", err)
		}

		err = k.reader.CommitMessages(ctx, kafkaMsg)
		if err != nil {
			return fmt.Errorf("commit kafka message: %w", err)
		}
	}
}

func (k kafkaConsumer) Close() error {
	err := k.reader.Close()
	if err != nil {
		return fmt.Errorf("close kafka consumer: %w", err)
	}

	return nil
}
