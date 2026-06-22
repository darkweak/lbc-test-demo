package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"leboncoin/pkg/services"
	"log"
	"net/http"
	"strconv"
)

type FizzBuzz struct {
	svcFizzBuzz   services.FizzBuzz
	svcStatistics services.Statistics
}

var _ Route = (*FizzBuzz)(nil)

func NewFizzBuzz(svcFizzBuzz services.FizzBuzz, svcStatistics services.Statistics) *FizzBuzz {
	return &FizzBuzz{
		svcFizzBuzz:   svcFizzBuzz,
		svcStatistics: svcStatistics,
	}
}

func (f *FizzBuzz) Register(server *http.ServeMux) {
	basepath := "/fizzbuzz"

	server.HandleFunc(basepath, f.queryParams)
	server.HandleFunc(basepath+"/{int1}/{int2}/{limit}/{str1}/{str2}", f.pathParams)
}

func (f *FizzBuzz) fizzbuzz(writer io.Writer, multiplyFirst, multiplySecond, limit int, fizzStr, buzzStr string) error {
	b, err := json.Marshal(f.svcFizzBuzz.Compute(multiplyFirst, multiplySecond, limit, fizzStr, buzzStr))
	if err != nil {
		return fmt.Errorf("failed to marshal fizzbuzz service: %w", err)
	}

	_, err = writer.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write fizzbuzz service: %w", err)
	}

	f.svcStatistics.Increment(fmt.Sprintf("%d-%d-%d-%s-%s", multiplyFirst, multiplySecond, limit, fizzStr, buzzStr))

	return nil
}

func (f *FizzBuzz) queryParams(writer http.ResponseWriter, request *http.Request) {
	multiplyFirst, err := strconv.Atoi(request.URL.Query().Get("int1"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	multiplySecond, err := strconv.Atoi(request.URL.Query().Get("int2"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	fizzStr := request.URL.Query().Get("str1")
	buzzStr := request.URL.Query().Get("str2")

	err = f.fizzbuzz(writer, multiplyFirst, multiplySecond, limit, fizzStr, buzzStr)
	if err != nil {
		log.Println("fizzbuzz failed:", err)

		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (f *FizzBuzz) pathParams(writer http.ResponseWriter, request *http.Request) {
	multiplyFirst, err := strconv.Atoi(request.PathValue("int1"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	multiplySecond, err := strconv.Atoi(request.PathValue("int2"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	limit, err := strconv.Atoi(request.PathValue("limit"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	fizzStr := request.URL.Query().Get("str1")
	buzzStr := request.URL.Query().Get("str2")

	err = f.fizzbuzz(writer, multiplyFirst, multiplySecond, limit, fizzStr, buzzStr)
	if err != nil {
		log.Println("fizzbuzz failed:", err)

		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}
