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

var _ Route = (*Statistics)(nil)

func NewStatistics(svcStatistics services.Statistics) *Statistics {
	return &Statistics{
		svcStatistics: svcStatistics,
	}
}

func (s *Statistics) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	statistic := s.svcStatistics.GetMostRecent()
	if statistic == nil {
		http.NotFound(writer, r)

		return
	}

	body, err := json.Marshal(statistic)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(body)
}

func (s *Statistics) Register(server *http.ServeMux) {
	server.Handle("/statistics", s)
}
