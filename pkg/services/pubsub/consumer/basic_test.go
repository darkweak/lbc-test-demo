package consumer_test

import (
	"errors"
	"leboncoin/pkg/services/pubsub"
	"leboncoin/pkg/services/pubsub/consumer"
	"testing"
	"time"
)

// Compile-time interface satisfaction check — NewBasicConsumer must return a pubsub.Consumer.
var _ pubsub.Consumer = consumer.NewBasicConsumer(nil)

var errCallbackFailed = errors.New("callback error")

func TestBasicConsumerInvokesCallbackForEachMessage(t *testing.T) {
	t.Parallel()

	queue := make(chan pubsub.Message, 3)
	queue <- pubsub.Message{Payload: []byte("a")}

	queue <- pubsub.Message{Payload: []byte("b")}

	queue <- pubsub.Message{Payload: []byte("c")}

	close(queue)

	c := consumer.NewBasicConsumer(queue)

	var received []string

	err := c.Consume(func(msg pubsub.Message) error {
		received = append(received, string(msg.Payload))

		return nil
	})
	if err != nil {
		t.Fatalf("Consume returned unexpected error: %v", err)
	}

	want := []string{"a", "b", "c"}
	if len(received) != len(want) {
		t.Fatalf("received %d messages, want %d", len(received), len(want))
	}

	for i, w := range want {
		if received[i] != w {
			t.Errorf("received[%d] = %q, want %q", i, received[i], w)
		}
	}
}

func TestBasicConsumerReturnsNilWhenChannelClosed(t *testing.T) {
	t.Parallel()

	queue := make(chan pubsub.Message)
	close(queue)

	c := consumer.NewBasicConsumer(queue)

	err := c.Consume(func(_ pubsub.Message) error {
		return nil
	})
	if err != nil {
		t.Errorf("Consume() = %v, want nil on closed channel", err)
	}
}

func TestBasicConsumerCallbackErrorReenqueuesAndReturnsWrappedError(t *testing.T) {
	t.Parallel()

	// Use a buffered channel of size 2: 1 for the initial message, 1 for the re-enqueued message.
	queue := make(chan pubsub.Message, 2)

	payload := []byte("fail-me")
	queue <- pubsub.Message{Payload: payload}

	c := consumer.NewBasicConsumer(queue)

	err := c.Consume(func(_ pubsub.Message) error {
		return errCallbackFailed
	})
	if err == nil {
		t.Fatal("Consume returned nil, want a wrapped error")
	}

	if !errors.Is(err, errCallbackFailed) {
		t.Errorf("Consume error = %v, want it to wrap %v", err, errCallbackFailed)
	}

	// The message must have been re-enqueued. Drain the channel with a timeout so
	// a broken implementation fails fast instead of blocking the test suite.
	select {
	case requeued := <-queue:
		if string(requeued.Payload) != string(payload) {
			t.Errorf("re-enqueued payload = %q, want %q", string(requeued.Payload), string(payload))
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for re-enqueued message — Consume may not have re-enqueued it")
	}
}

func TestBasicConsumerCloseReturnsNil(t *testing.T) {
	t.Parallel()

	queue := make(chan pubsub.Message)
	c := consumer.NewBasicConsumer(queue)

	err := c.Close()
	if err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}
}
