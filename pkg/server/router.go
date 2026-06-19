package server

import (
	"leboncoin/pkg/server/routes"
	"net/http"
)

func NewRouter(server *http.ServeMux, appRoutes ...routes.Route) {
	for _, route := range appRoutes {
		route.Register(server)
	}
}
