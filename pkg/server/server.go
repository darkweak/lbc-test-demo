package server

import (
	"leboncoin/pkg/configuration"
	"net/http"
)

func NewServer(cfg configuration.Configuration) *http.Server {
	server := new(http.Server)
	server.Addr = cfg.Address

	server.Handler = http.NewServeMux()

	return server
}
