//go:build integration

package consumer_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	kafkacontainer "github.com/testcontainers/testcontainers-go/modules/kafka"

	"leboncoin/pkg/services/pubsub"
	"leboncoin/pkg/services/pubsub/consumer"

	kafka "github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
)

const (
	consumerTestTopic     = "consumer-integration-test"
	consumeTimeout        = 30 * time.Second
	containerStartTimeout = 2 * time.Minute
)

// errStop is returned from the callback after all expected messages are received,
// causing Consume to break out of its infinite loop.
var errStop = errors.New("stop consuming")

// startKafkaForConsumer starts a single-node KRaft Kafka container and returns its broker
// address. The container is registered for cleanup via t.Cleanup.
func startKafkaForConsumer(t *testing.T) []string {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), containerStartTimeout)
	defer cancel()

	kafkaCtr, err := kafkacontainer.Run(ctx, "confluentinc/cp-kafka:7.8.0",
		kafkacontainer.WithClusterID("test-cluster-consumer"),
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

// createTopicForConsumer creates a Kafka topic via a raw kafka-go connection.
func createTopicForConsumer(t *testing.T, brokers []string, topic string) {
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

// writeJSONMessages writes JSON-encoded pubsub.Message values via a raw kafka-go writer.
// The consumer JSON-unmarshals each message value into pubsub.Message, so raw bytes
// are not compatible — the test must write properly encoded JSON.
func writeJSONMessages(t *testing.T, brokers []string, topic string, payloads [][]byte) {
	t.Helper()

	writer := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}
	defer func() {
		err := writer.Close()
		if err != nil {
			t.Logf("close raw writer: %v", err)
		}
	}()

	msgs := make([]kafka.Message, 0, len(payloads))

	for _, payload := range payloads {
		encoded, err := json.Marshal(pubsub.Message{Payload: payload})
		if err != nil {
			t.Fatalf("marshal pubsub.Message: %v", err)
		}

		msgs = append(msgs, kafka.Message{Value: encoded})
	}

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()

	err := writer.WriteMessages(ctx, msgs...)
	if err != nil {
		t.Fatalf("write messages: %v", err)
	}
}

// collectMessages drains the received channel until wantCount messages are collected or
// the timeout fires. It returns the collected messages.
func collectMessages(t *testing.T, received <-chan pubsub.Message, wantCount int) []pubsub.Message {
	t.Helper()

	timer := time.NewTimer(consumeTimeout)
	defer timer.Stop()

	got := make([]pubsub.Message, 0, wantCount)

	for len(got) < wantCount {
		select {
		case msg := <-received:
			got = append(got, msg)
		case <-timer.C:
			t.Fatalf("timed out after %s waiting for messages; received %d of %d", consumeTimeout, len(got), wantCount)
		}
	}

	return got
}

func TestKafkaConsumerConsumeDeliverstMessagesToCallback(t *testing.T) {
	t.Parallel()

	brokers := startKafkaForConsumer(t)
	createTopicForConsumer(t, brokers, consumerTestTopic)

	payloads := [][]byte{
		[]byte("first"),
		[]byte("second"),
		[]byte("third"),
	}

	writeJSONMessages(t, brokers, consumerTestTopic, payloads)

	kafkaConsumer := consumer.NewKafkaConsumer(brokers, consumerTestTopic)

	t.Cleanup(func() {
		err := kafkaConsumer.Close()
		if err != nil {
			t.Logf("close consumer: %v", err)
		}
	})

	received := make(chan pubsub.Message, len(payloads))
	consumeErr := make(chan error, 1)

	go func() {
		seen := 0

		consumeErr <- kafkaConsumer.Consume(func(msg pubsub.Message) error {
			received <- msg

			seen++

			// Return errStop after all expected messages are collected.
			// Consume wraps this and returns it, breaking the loop.
			if seen == len(payloads) {
				return errStop
			}

			return nil
		})
	}()

	got := collectMessages(t, received, len(payloads))

	// Wait for Consume to exit and verify it returned the sentinel wrapped error.
	select {
	case err := <-consumeErr:
		if !errors.Is(err, errStop) {
			t.Errorf("Consume returned error = %v, want it to wrap errStop", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for Consume to exit after sentinel error")
	}

	for i, want := range payloads {
		if string(got[i].Payload) != string(want) {
			t.Errorf("message[%d].Payload = %q, want %q", i, string(got[i].Payload), string(want))
		}
	}
}

func TestKafkaConsumerCloseReturnsNoError(t *testing.T) {
	t.Parallel()

	brokers := startKafkaForConsumer(t)
	createTopicForConsumer(t, brokers, consumerTestTopic+"-close")

	kafkaConsumer := consumer.NewKafkaConsumer(brokers, consumerTestTopic+"-close")

	err := kafkaConsumer.Close()
	if err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}
}
