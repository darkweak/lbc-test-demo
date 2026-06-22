package consumer

import (
	"fmt"
	"leboncoin/pkg/services/pubsub"
)

type basicConsumer struct {
	queue chan pubsub.Message
}

var _ pubsub.Consumer = (*basicConsumer)(nil)

func NewBasicConsumer(queue chan pubsub.Message) *basicConsumer {
	return &basicConsumer{
		queue: queue,
	}
}

func (k basicConsumer) Consume(callback func(message pubsub.Message) error) error {
	for {
		msg, ok := <-k.queue
		if !ok {
			return nil
		}

		err := callback(msg)
		if err != nil {
			k.queue <- msg

			return fmt.Errorf("error while processing message: %w", err)
		}
	}
}

func (k basicConsumer) Close() error {
	return nil
}
