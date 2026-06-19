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

func NewFizzBuzz(svcFizzBuzz services.FizzBuzz, svcStatistics services.Statistics) Route {
	return &FizzBuzz{
		svcFizzBuzz:   svcFizzBuzz,
		svcStatistics: svcStatistics,
	}
}

func (f *FizzBuzz) fizzbuzz(w io.Writer, multiplyFirst, multiplySecond, limit int, fizzStr, buzzStr string) error {
	b, err := json.Marshal(f.svcFizzBuzz.Compute(multiplyFirst, multiplySecond, limit, fizzStr, buzzStr))
	if err != nil {
		return fmt.Errorf("failed to marshal fizzbuzz service: %w", err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write fizzbuzz service: %w", err)
	}

	f.svcStatistics.Increment(fmt.Sprintf("%d-%d-%d-%s-%s", multiplyFirst, multiplySecond, limit, fizzStr, buzzStr))

	return nil
}

func (f *FizzBuzz) queryParams(w http.ResponseWriter, r *http.Request) {
	multiplyFirst, err := strconv.Atoi(r.URL.Query().Get("int1"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	multiplySecond, err := strconv.Atoi(r.URL.Query().Get("int2"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	fizzStr := r.URL.Query().Get("str1")
	buzzStr := r.URL.Query().Get("str2")

	err = f.fizzbuzz(w, multiplyFirst, multiplySecond, limit, fizzStr, buzzStr)
	if err != nil {
		log.Println("fizzbuzz failed:", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (f *FizzBuzz) pathParams(w http.ResponseWriter, r *http.Request) {
	multiplyFirst, err := strconv.Atoi(r.PathValue("int1"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	multiplySecond, err := strconv.Atoi(r.PathValue("int2"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	limit, err := strconv.Atoi(r.PathValue("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	fizzStr := r.URL.Query().Get("str1")
	buzzStr := r.URL.Query().Get("str2")

	err = f.fizzbuzz(w, multiplyFirst, multiplySecond, limit, fizzStr, buzzStr)
	if err != nil {
		log.Println("fizzbuzz failed:", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (f *FizzBuzz) Register(server *http.ServeMux) {
	basepath := "/fizzbuzz"

	server.HandleFunc(basepath, f.queryParams)
	server.HandleFunc(basepath+"/{int1}/{int2}/{limit}/{str1}/{str2}", f.pathParams)
}
