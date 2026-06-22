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

	svcStatistics := services.NewStatistics()

	mux, ok := srv.Handler.(*http.ServeMux)
	if !ok {
		_ = srv.Close()

		log.Fatal("http server does not implement http.ServeMux")
	}

	defer func() {
		_ = srv.Close()
	}()

	server.NewRouter(
		mux,
		routes.NewFizzBuzz(services.NewFizzBuzz(), svcStatistics),
		routes.NewStatistics(svcStatistics),
	)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
