package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"leboncoin/pkg/server/middleware"
	"leboncoin/pkg/services/pubsub"

	kafka "github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

var _ pubsub.Producer = (*kafkaProducer)(nil)

func NewKafkaProducer(hosts []string, topic string) *kafkaProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(hosts...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (k kafkaProducer) Produce(ctx context.Context, message []byte) error {
	value, err := json.Marshal(pubsub.Message{Payload: message})
	if err != nil {
		return fmt.Errorf("marshal kafka message: %w", err)
	}

	correlationID, _ := ctx.Value(middleware.CorrelationIDCtxKey).(string)

	err = k.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
		Headers: []kafka.Header{{
			Key:   "Correlation-ID",
			Value: []byte(correlationID),
		}},
	})
	if err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}

	return nil
}

func (k kafkaProducer) Close() error {
	err := k.writer.Close()
	if err != nil {
		return fmt.Errorf("close kafka producer: %w", err)
	}

	return nil
}
