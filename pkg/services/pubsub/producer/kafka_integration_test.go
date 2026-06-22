//go:build integration

package producer_test

import (
	"context"
	"testing"
	"time"

	kafkacontainer "github.com/testcontainers/testcontainers-go/modules/kafka"

	"leboncoin/pkg/services/pubsub/producer"

	kafka "github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
)

const (
	producerTestTopic     = "producer-integration-test"
	producerReadTimeout   = 30 * time.Second
	containerStartTimeout = 2 * time.Minute
)

// startKafkaForProducer starts a single-node KRaft Kafka container and returns its broker
// address. The container is registered for cleanup via t.Cleanup.
func startKafkaForProducer(t *testing.T) []string {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), containerStartTimeout)
	defer cancel()

	kafkaCtr, err := kafkacontainer.Run(ctx, "confluentinc/cp-kafka:7.8.0",
		kafkacontainer.WithClusterID("test-cluster-producer"),
	)
	testcontainers.CleanupContainer(t, kafkaCtr)

	if err != nil {
		t.Skipf("skipping integration test: could not start Kafka container: %v", err)
	}

	brokers, err := kafkaCtr.Brokers(t.Context())
	if err != nil {
		t.Fatalf("get kafka brokers: %v", err)
	}

	return brokers
}

// createTopicForProducer creates a Kafka topic via a raw kafka-go connection.
func createTopicForProducer(t *testing.T, brokers []string, topic string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		t.Fatalf("dial kafka: %v", err)
	}

	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			t.Logf("close kafka conn: %v", closeErr)
		}
	}()

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		t.Fatalf("create topic %q: %v", topic, err)
	}
}

// readOneMessage reads a single message from a raw kafka-go reader on the given topic
// and returns it. The reader is registered for cleanup via t.Cleanup.
func readOneMessage(t *testing.T, brokers []string, topic string) kafka.Message {
	t.Helper()

	rawReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  1 << 20,
	})

	t.Cleanup(func() {
		err := rawReader.Close()
		if err != nil {
			t.Logf("close raw reader: %v", err)
		}
	})

	ctx, cancel := context.WithTimeout(t.Context(), producerReadTimeout)
	defer cancel()

	msg, err := rawReader.ReadMessage(ctx)
	if err != nil {
		t.Fatalf("read message: %v", err)
	}

	return msg
}

func TestKafkaProducerProduceWritesMessageWithCorrelationIDHeader(t *testing.T) {
	t.Parallel()

	brokers := startKafkaForProducer(t)
	createTopicForProducer(t, brokers, producerTestTopic)

	prod := producer.NewKafkaProducer(brokers, producerTestTopic)

	t.Cleanup(func() {
		err := prod.Close()
		if err != nil {
			t.Logf("close producer: %v", err)
		}
	})

	payload := []byte(`hello integration`)

	err := prod.Produce(payload)
	if err != nil {
		t.Fatalf("Produce: %v", err)
	}

	// Read the message back with a raw kafka-go reader — independent of the package consumer.
	msg := readOneMessage(t, brokers, producerTestTopic)

	// Assert the value matches what was produced.
	if string(msg.Value) != string(payload) {
		t.Errorf("message Value = %q, want %q", string(msg.Value), string(payload))
	}

	// Assert the Correlation-ID header is present and non-empty.
	var correlationID string

	for _, header := range msg.Headers {
		if header.Key == "Correlation-ID" {
			correlationID = string(header.Value)

			break
		}
	}

	if correlationID == "" {
		t.Error("expected Correlation-ID header to be present and non-empty, got none")
	}
}

func TestKafkaProducerCloseReturnsNoError(t *testing.T) {
	t.Parallel()

	brokers := startKafkaForProducer(t)
	createTopicForProducer(t, brokers, producerTestTopic+"-close")

	prod := producer.NewKafkaProducer(brokers, producerTestTopic+"-close")

	err := prod.Close()
	if err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}
}
