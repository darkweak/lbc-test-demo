package pubsub

import "io"

type Message struct {
	Payload []byte `json:"payload"`
}

type Producer interface {
	io.Closer

	Produce(message []byte) error
}

type Consumer interface {
	io.Closer

	Consume(callback func(msg Message) error) error
}

type PubSub interface {
	Producer
	Consumer
}
