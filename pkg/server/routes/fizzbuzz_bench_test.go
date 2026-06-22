package routes_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"leboncoin/pkg/services"
)

const (
	benchLimitSmall  = "small/limit=15"
	benchLimitMedium = "medium/limit=1000"
)

func BenchmarkFizzBuzzHandlerQueryParams(b *testing.B) {
	benchCases := []struct {
		name string
		url  string
		stub []string
	}{
		{
			name: benchLimitSmall,
			url:  "/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz",
			stub: makeStubResult(15),
		},
		{
			name: benchLimitMedium,
			url:  "/fizzbuzz?int1=3&int2=5&limit=1000&str1=fizz&str2=buzz",
			stub: makeStubResult(1_000),
		},
		{
			name: "large/limit=100000",
			url:  "/fizzbuzz?int1=3&int2=5&limit=100000&str1=fizz&str2=buzz",
			stub: makeStubResult(100_000),
		},
	}

	for _, benchCase := range benchCases {
		b.Run(benchCase.name, func(b *testing.B) {
			mux := newFizzBuzzMux(
				&stubFizzBuzz{result: benchCase.stub},
				&stubStatistics{incrementedKeys: nil, mostRecent: nil},
			)

			req := httptest.NewRequestWithContext(
				context.Background(), http.MethodGet, benchCase.url, nil,
			)

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				rec := httptest.NewRecorder()
				mux.ServeHTTP(rec, req)
			}
		})
	}
}

func BenchmarkFizzBuzzHandlerPathParams(b *testing.B) {
	benchCases := []struct {
		name string
		url  string
		stub []string
	}{
		{
			name: benchLimitSmall,
			url:  "/fizzbuzz/3/5/15/fizz/buzz?str1=fizz&str2=buzz",
			stub: makeStubResult(15),
		},
		{
			name: benchLimitMedium,
			url:  "/fizzbuzz/3/5/1000/fizz/buzz?str1=fizz&str2=buzz",
			stub: makeStubResult(1_000),
		},
	}

	for _, benchCase := range benchCases {
		b.Run(benchCase.name, func(b *testing.B) {
			mux := newFizzBuzzMux(
				&stubFizzBuzz{result: benchCase.stub},
				&stubStatistics{incrementedKeys: nil, mostRecent: nil},
			)

			req := httptest.NewRequestWithContext(
				context.Background(), http.MethodGet, benchCase.url, nil,
			)

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				rec := httptest.NewRecorder()
				mux.ServeHTTP(rec, req)
			}
		})
	}
}

func makeStubResult(count int) []string {
	out := make([]string, count)
	for idx := range out {
		out[idx] = "fizz"
	}

	return out
}

func BenchmarkFizzBuzzHandlerBadRequest(b *testing.B) {
	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubStatistics{incrementedKeys: nil, mostRecent: nil},
	)

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/fizzbuzz?int2=5&limit=15&str1=fizz&str2=buzz",
		nil,
	)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
	}
}

func BenchmarkFizzBuzzHandlerQueryParamsRealService(b *testing.B) {
	benchCases := []struct {
		name string
		url  string
	}{
		{benchLimitSmall, "/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz"},
		{benchLimitMedium, "/fizzbuzz?int1=3&int2=5&limit=1000&str1=fizz&str2=buzz"},
	}

	for _, benchCase := range benchCases {
		b.Run(benchCase.name, func(b *testing.B) {
			mux := newFizzBuzzMux(
				services.NewFizzBuzz(),
				&stubStatistics{incrementedKeys: nil, mostRecent: nil},
			)

			req := httptest.NewRequestWithContext(
				context.Background(), http.MethodGet, benchCase.url, nil,
			)

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				rec := httptest.NewRecorder()
				mux.ServeHTTP(rec, req)
			}
		})
	}
}
