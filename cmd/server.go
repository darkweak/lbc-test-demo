package main

import (
	"leboncoin/pkg/configuration"
	"leboncoin/pkg/server"
	"leboncoin/pkg/server/routes"
	"leboncoin/pkg/services/fizzbuzz"
	"leboncoin/pkg/services/pubsub"
	"leboncoin/pkg/services/pubsub/consumer"
	"leboncoin/pkg/services/pubsub/producer"
	"leboncoin/pkg/services/statistics"
	"log"
	"net/http"
)

func getPubsub(cfg configuration.Configuration) (fizzBuzzConsumer pubsub.Consumer, fizzBuzzProducer pubsub.Producer) {
	if cfg.FTKafka {
		topic := "my-topic"

		fizzBuzzProducer = producer.NewKafkaProducer(cfg.KafkaAddresses, topic)
		fizzBuzzConsumer = consumer.NewKafkaConsumer(cfg.KafkaAddresses, topic)
	} else {
		fizzbuzzQueue := make(chan pubsub.Message)

		fizzBuzzProducer = producer.NewBasicProducer(fizzbuzzQueue)
		fizzBuzzConsumer = consumer.NewBasicConsumer(fizzbuzzQueue)
	}

	return
}

func main() {
	cfg, err := configuration.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(cfg)

	mux, ok := srv.Handler.(*http.ServeMux)
	if !ok {
		_ = srv.Close()

		log.Fatal("http server does not implement http.ServeMux")
	}

	svcStatistics := statistics.NewStatistics()
	fizzBuzzConsumer, fizzBuzzProducer := getPubsub(cfg)

	log.Printf("Running pubsub consumer %T and producer %T", fizzBuzzConsumer, fizzBuzzProducer)

	defer func() {
		_ = srv.Close()
		_ = fizzBuzzProducer.Close()
		_ = fizzBuzzConsumer.Close()
	}()

	go func() {
		_ = fizzBuzzConsumer.Consume(func(message pubsub.Message) error {
			svcStatistics.Increment(string(message.Payload))

			return nil
		})
	}()

	fizzBuzzRoute := routes.NewFizzBuzz(fizzbuzz.NewFizzBuzz(), fizzBuzzProducer)
	server.NewRouter(mux, fizzBuzzRoute, routes.NewStatistics(svcStatistics))

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
