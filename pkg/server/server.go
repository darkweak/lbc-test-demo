package server

import (
	"leboncoin/pkg/configuration"
	"net/http"
	"slices"
)

func NewServer(cfg configuration.Configuration, middlewares ...func(handler http.Handler) http.Handler) *http.Server {
	server := new(http.Server)
	server.Addr = cfg.Address
	server.Handler = http.NewServeMux()

	for _, middleware := range slices.Backward(middlewares) {
		server.Handler = middleware(server.Handler)
	}

	return server
}
