package main

import (
	"leboncoin/pkg/configuration"
	"leboncoin/pkg/server"
	"leboncoin/pkg/server/routes"
	"leboncoin/pkg/services"
	"log"
	"net/http"
)

func main() {
	cfg, err := configuration.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(cfg)
	defer func() {
		_ = srv.Close()
	}()

	svcStatistics := services.NewStatistics()

	server.NewRouter(
		srv.Handler.(*http.ServeMux),
		routes.NewFizzBuzz(services.NewFizzBuzz(), svcStatistics),
		routes.NewStatistics(svcStatistics),
	)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
