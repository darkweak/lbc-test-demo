package routes

import (
	"encoding/json"
	"leboncoin/pkg/services"
	"log"
	"net/http"
)

type Statistics struct {
	svcStatistics services.Statistics
}

func NewStatistics(svcStatistics services.Statistics) Route {
	return &Statistics{
		svcStatistics: svcStatistics,
	}
}

func (s *Statistics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statistic := s.svcStatistics.GetMostRecent()
	if statistic == nil {
		http.NotFound(w, r)

		return
	}

	body, err := json.Marshal(statistic)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (s *Statistics) Register(server *http.ServeMux) {
	server.Handle("/statistics", s)
}
