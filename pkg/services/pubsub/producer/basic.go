package producer

import (
	"leboncoin/pkg/services/pubsub"
)

type basicProducer struct {
	queue chan pubsub.Message
}

var _ pubsub.Producer = (*basicProducer)(nil)

func NewBasicProducer(queue chan pubsub.Message) *basicProducer {
	return &basicProducer{
		queue: queue,
	}
}

func (k basicProducer) Produce(message []byte) error {
	k.queue <- pubsub.Message{Payload: message}

	return nil
}

func (k basicProducer) Close() error {
	return nil
}
