package producer_test

import (
	"leboncoin/pkg/services/pubsub"
	"leboncoin/pkg/services/pubsub/producer"
	"testing"
)

// Compile-time interface satisfaction check — NewBasicProducer must return a pubsub.Producer.
var _ pubsub.Producer = producer.NewBasicProducer(nil)

func TestBasicProducerSendsMessageToChannel(t *testing.T) {
	t.Parallel()

	queue := make(chan pubsub.Message, 1)
	p := producer.NewBasicProducer(queue)

	payload := []byte("hello")

	err := p.Produce(t.Context(), payload)
	if err != nil {
		t.Fatalf("Produce returned unexpected error: %v", err)
	}

	select {
	case msg := <-queue:
		if string(msg.Payload) != string(payload) {
			t.Errorf("message payload = %q, want %q", string(msg.Payload), string(payload))
		}
	default:
		t.Fatal("channel is empty after Produce — message was not sent")
	}
}

func TestBasicProducerCloseReturnsNil(t *testing.T) {
	t.Parallel()

	queue := make(chan pubsub.Message, 1)
	p := producer.NewBasicProducer(queue)

	err := p.Close()
	if err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}
}
