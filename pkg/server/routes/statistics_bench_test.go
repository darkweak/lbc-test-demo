package routes_test

import (
	"context"
	"leboncoin/pkg/services/statistics"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const benchStatKey = "3-5-15-fizz-buzz"

func BenchmarkStatisticsHandlerNoData(b *testing.B) {
	mux := newStatisticsMux(&stubStatistics{incrementedKeys: nil, mostRecent: nil})
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/statistics", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
	}
}

func BenchmarkStatisticsHandlerWithData(b *testing.B) {
	stat := &statistics.Statistic{
		LastCall: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Hit:      42,
		Key:      benchStatKey,
	}

	mux := newStatisticsMux(&stubStatistics{incrementedKeys: nil, mostRecent: stat})
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/statistics", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
	}
}

func BenchmarkStatisticsHandlerWithDataParallel(b *testing.B) {
	stat := &statistics.Statistic{
		LastCall: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Hit:      42,
		Key:      benchStatKey,
	}

	mux := newStatisticsMux(&stubStatistics{incrementedKeys: nil, mostRecent: stat})
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/statistics", nil)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
		}
	})
}

func BenchmarkStatisticsHandlerRealService(b *testing.B) {
	svc := statistics.NewStatistics()
	svc.Increment(benchStatKey)
	svc.Increment(benchStatKey)

	mux := newStatisticsMux(svc)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/statistics", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
	}
}
